package raedian

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/modbus"
)

// Modbus Register Addresses
const (
	raedianRegSerial               = 0x8000
	raedianRegFirmwareVersion      = 0x8004
	raedianRegMaxRatedCurrent      = 0x8006
	raedianRegErrorCode            = 0x8008
	raedianRegSocketLockState      = 0x800A
	raedianRegChargingState        = 0x800C
	raedianRegCurrentChargingLimit = 0x800E
	raedianRegChargingCurrentL1    = 0x8010
	raedianRegChargingCurrentL2    = 0x8012
	raedianRegChargingCurrentL3    = 0x8014
	raedianRegVoltageL1            = 0x8016
	raedianRegVoltageL2            = 0x8018
	raedianRegVoltageL3            = 0x801A
	raedianRegActivePower          = 0x801C
	raedianRegEnergyDelivered      = 0x801E
	raedianRegSetChargingCurrent   = 0x8100
	raedianRegSetChargingPhase     = 0x8102
	raedianRegStartStopSession     = 0x8105
)

// Charger is an api.Charger implementation
type Charger struct {
	conn *modbus.Connection
	log  *util.Logger
}

// NewRaedianFromConfig creates a new Raedian charger
func NewRaedianFromConfig(ctx context.Context, other map[string]interface{}) (api.Charger, error) {
	cc := modbus.Settings{
		ID: 1,
	}
	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}
	return NewRaedian(ctx, cc.URI, cc.Device, cc.Comset, cc.Baudrate, cc.Protocol(), cc.ID)
}

// NewRaedian creates a new charger
func NewRaedian(ctx context.Context, uri, device, comset string, baudrate int, proto modbus.Protocol, slaveID uint8) (*Charger, error) {
	conn, err := modbus.NewConnection(ctx, uri, device, comset, baudrate, proto, slaveID)
	if err != nil {
		return nil, err
	}
	log := util.NewLogger("raedian")
	conn.Logger(log.TRACE)

	wb := &Charger{
		conn: conn,
		log:  log,
	}

	return wb, err
}

// Status implements the api.Charger interface
func (c *Charger) Status() (api.ChargeStatus, error) {
	b, err := c.conn.ReadHoldingRegisters(raedianRegChargingState, 1)
	if err != nil {
		return api.StatusNone, err
	}
	status := binary.BigEndian.Uint16(b)

	switch status {
	case 0:
		return api.StatusA, nil
	case 1, 2:
		return api.StatusB, nil
	case 3, 4:
		return api.StatusC, nil
	default:
		return api.StatusNone, fmt.Errorf("unknown status: %d", status)
	}
}

// Enabled implements the api.Charger interface
func (c *Charger) Enabled() (bool, error) {
	b, err := c.conn.ReadHoldingRegisters(raedianRegCurrentChargingLimit, 1)
	if err != nil {
		return false, err
	}
	return binary.BigEndian.Uint16(b) >= 6000, nil
}

// Enable implements the api.Charger interface
func (c *Charger) Enable(enable bool) error {
	var value uint16
	if !enable {
		value = 0x01
	}
	_, err := c.conn.WriteSingleRegister(raedianRegStartStopSession, value)
	if err != nil {
		return err
	}
	return nil
}

// MaxCurrent implements the api.Charger interface
func (c *Charger) MaxCurrent(current int64) error {
	if current < 6 {
		current = 0
	}
	_, err := c.conn.WriteSingleRegister(raedianRegSetChargingCurrent, uint16(current * 1000))
	return err
}

var _ api.Meter = (*Charger)(nil)

// CurrentPower implements the api.Meter interface
func (c *Charger) CurrentPower() (float64, error) {
	b, err := c.conn.ReadHoldingRegisters(raedianRegActivePower, 1)
	if err != nil {
		return 0, err
	}
	return float64(binary.BigEndian.Uint16(b)), nil
}

var _ api.MeterEnergy = (*Charger)(nil)

// TotalEnergy implements the api.MeterEnergy interface
func (c *Charger) TotalEnergy() (float64, error) {
	b, err := c.conn.ReadHoldingRegisters(raedianRegEnergyDelivered, 1)
	if err != nil {
		return 0, err
	}
	return float64(binary.BigEndian.Uint16(b)) / 1000, nil
}

var _ api.PhaseCurrents = (*Charger)(nil)

// Currents implements the api.PhaseCurrents interface
func (c *Charger) Currents() (float64, float64, float64, error) {
	b, err := c.conn.ReadHoldingRegisters(raedianRegChargingCurrentL1, 3)
	if err != nil {
		return 0, 0, 0, err
	}
	l1 := float64(binary.BigEndian.Uint16(b[0:2])) / 1000
	l2 := float64(binary.BigEndian.Uint16(b[2:4])) / 1000
	l3 := float64(binary.BigEndian.Uint16(b[4:6])) / 1000
	return l1, l2, l3, nil
}

var _ api.PhaseVoltages = (*Charger)(nil)

// Voltages implements the api.PhaseVoltages interface
func (c *Charger) Voltages() (float64, float64, float64, error) {
	b, err := c.conn.ReadHoldingRegisters(raedianRegVoltageL1, 3)
	if err != nil {
		return 0, 0, 0, err
	}
	l1 := float64(binary.BigEndian.Uint16(b[0:2])) / 10
	l2 := float64(binary.BigEndian.Uint16(b[2:4])) / 10
	l3 := float64(binary.BigEndian.Uint16(b[4:6])) / 10
	return l1, l2, l3, nil
}

var _ api.Diagnosis = (*Charger)(nil)

// Diagnose implements the api.Diagnosis interface
func (c *Charger) Diagnose() {
	if b, err := c.conn.ReadHoldingRegisters(raedianRegSerial, 2); err == nil {
		fmt.Printf("\tSerial: %d\n", binary.BigEndian.Uint32(b))
	}

	if b, err := c.conn.ReadHoldingRegisters(raedianRegFirmwareVersion, 1); err == nil {
		fmt.Printf("\tFirmware Version: %d\n", binary.BigEndian.Uint16(b))
	} else {
		fmt.Printf("\tFirmware Version: ERROR - %v\n", err)
	}

	if b, err := c.conn.ReadHoldingRegisters(raedianRegMaxRatedCurrent, 1); err == nil {
		fmt.Printf("\tMax Rated Current: %.3f A\n", float64(binary.BigEndian.Uint16(b))/1000)
	} else {
		fmt.Printf("\tMax Rated Current: ERROR - %v\n", err)
	}

	if b, err := c.conn.ReadHoldingRegisters(raedianRegErrorCode, 1); err == nil {
		fmt.Printf("\tError Code: %d\n", binary.BigEndian.Uint16(b))
	} else {
		fmt.Printf("\tError Code: ERROR - %v\n", err)
	}

	if b, err := c.conn.ReadHoldingRegisters(raedianRegSocketLockState, 1); err == nil {
		fmt.Printf("\tSocket Lock State: %d\n", binary.BigEndian.Uint16(b))
	} else {
		fmt.Printf("\tSocket Lock State: ERROR - %v\n", err)
	}
}
