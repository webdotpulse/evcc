<template>
	<div class="visual-energyflow" ref="container">
		<div class="node-container">
			<!-- Top Row -->
			<div class="row top-row">
				<div class="node grid-node" ref="gridNode" :class="{ active: gridImport > 0 || gridExport > 0 }">
					<div class="icon"><shopicon-regular-powersupply></shopicon-regular-powersupply></div>
					<div class="label">Grid</div>
					<div class="value">
						<span v-if="gridImport > 0">{{ fmtW(gridImport, powerUnit) }} In</span>
						<span v-else-if="gridExport > 0">{{ fmtW(gridExport, powerUnit) }} Out</span>
						<span v-else>0 W</span>
					</div>
				</div>
				<div class="node pv-node" ref="pvNode" :class="{ active: pvProduction > 0 }">
					<div class="icon"><shopicon-regular-sun></shopicon-regular-sun></div>
					<div class="label">Solar</div>
					<div class="value">{{ fmtW(pvProduction, powerUnit) }}</div>
				</div>
			</div>

			<!-- Center Row -->
			<div class="row center-row">
				<div class="node home-node active" ref="homeNode">
					<div class="icon"><shopicon-regular-home></shopicon-regular-home></div>
					<div class="label">Home</div>
					<div class="value">{{ fmtW(homePower, powerUnit) }}</div>
				</div>
			</div>

			<!-- Bottom Row -->
			<div class="row bottom-row">
				<div class="node battery-node" ref="batteryNode" v-if="batteryConfigured" :class="{ active: batteryPower !== 0 }">
					<div class="icon"><BatteryIcon :soc="batterySoc" /></div>
					<div class="label">Battery ({{ fmtPercentage(batterySoc) }})</div>
					<div class="value">
						<span v-if="batteryDischarge > 0">{{ fmtW(batteryDischarge, powerUnit) }} Out</span>
						<span v-else-if="batteryCharge > 0">{{ fmtW(batteryCharge, powerUnit) }} In</span>
						<span v-else>0 W</span>
					</div>
				</div>
				<!-- Empty div to keep alignment if no battery -->
				<div v-else style="width: 100px;"></div>

				<div class="node vehicle-node" ref="vehicleNode" v-if="loadpoints.length > 0" :class="{ active: loadpointsPower > 0 }">
					<div class="icon"><VehicleIcon :names="vehicleIcons" /></div>
					<div class="label">Vehicles</div>
					<div class="value">{{ fmtW(loadpointsPower, powerUnit) }}</div>
				</div>
				<!-- Empty div to keep alignment if no loadpoints -->
				<div v-else style="width: 100px;"></div>
			</div>
		</div>

		<svg class="flow-lines" width="100%" height="100%" v-if="paths.gridToHome">
			<defs>
				<marker id="arrow-grid-in" viewBox="0 0 10 10" refX="5" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
					<path d="M 0 0 L 10 5 L 0 10 z" fill="var(--evcc-grid)" />
				</marker>
				<marker id="arrow-grid-out" viewBox="0 0 10 10" refX="5" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
					<path d="M 0 0 L 10 5 L 0 10 z" fill="var(--evcc-export)" />
				</marker>
				<marker id="arrow-pv" viewBox="0 0 10 10" refX="5" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
					<path d="M 0 0 L 10 5 L 0 10 z" fill="var(--evcc-pv)" />
				</marker>
				<marker id="arrow-battery-charge" viewBox="0 0 10 10" refX="5" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
					<path d="M 0 0 L 10 5 L 0 10 z" fill="var(--evcc-battery)" />
				</marker>
				<marker id="arrow-battery-discharge" viewBox="0 0 10 10" refX="5" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
					<path d="M 0 0 L 10 5 L 0 10 z" fill="var(--evcc-battery)" />
				</marker>
				<marker id="arrow-vehicle" viewBox="0 0 10 10" refX="5" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
					<path d="M 0 0 L 10 5 L 0 10 z" fill="var(--evcc-export)" />
				</marker>
			</defs>

			<!-- Grid <-> Home -->
			<g v-if="paths.gridToHome">
				<path :d="paths.gridToHome"
					  class="flow-path"
					  :class="{ 'grid-in': gridImport > 0, 'grid-out': gridExport > 0 }"
					  :marker-end="gridImport > 0 ? 'url(#arrow-grid-in)' : (gridExport > 0 ? 'url(#arrow-grid-out)' : '')" />
				<!-- Animated dots for flow -->
				<circle v-if="gridImport > 0" r="6" class="dot grid-dot-in">
					<animateMotion :path="paths.gridToHome" dur="2s" repeatCount="indefinite" />
				</circle>
				<circle v-if="gridExport > 0" r="6" class="dot grid-dot-out">
					<animateMotion :path="paths.homeToGrid" dur="2s" repeatCount="indefinite" />
				</circle>
			</g>

			<!-- PV -> Home -->
			<g v-if="paths.pvToHome">
				<path :d="paths.pvToHome"
					  class="flow-path"
					  :class="{ 'pv-in': pvProduction > 0 }"
					  :marker-end="pvProduction > 0 ? 'url(#arrow-pv)' : ''" />
				<circle v-if="pvProduction > 0" r="6" class="dot pv-dot">
					<animateMotion :path="paths.pvToHome" dur="2s" repeatCount="indefinite" />
				</circle>
			</g>

			<!-- Battery <-> Home -->
			<g v-if="batteryConfigured && paths.homeToBattery">
				<path :d="paths.homeToBattery"
					  class="flow-path"
					  :class="{ 'battery-charge': batteryCharge > 0, 'battery-discharge': batteryDischarge > 0 }"
					  :marker-end="batteryCharge > 0 ? 'url(#arrow-battery-charge)' : (batteryDischarge > 0 ? 'url(#arrow-battery-discharge)' : '')" />
				<circle v-if="batteryCharge > 0" r="6" class="dot battery-dot">
					<animateMotion :path="paths.homeToBattery" dur="2s" repeatCount="indefinite" />
				</circle>
				<circle v-if="batteryDischarge > 0" r="6" class="dot battery-dot">
					<animateMotion :path="paths.batteryToHome" dur="2s" repeatCount="indefinite" />
				</circle>
			</g>

			<!-- Home -> Vehicle -->
			<g v-if="loadpoints.length > 0 && paths.homeToVehicle">
				<path :d="paths.homeToVehicle"
					  class="flow-path"
					  :class="{ 'vehicle-in': loadpointsPower > 0 }"
					  :marker-end="loadpointsPower > 0 ? 'url(#arrow-vehicle)' : ''" />
				<circle v-if="loadpointsPower > 0" r="6" class="dot vehicle-dot">
					<animateMotion :path="paths.homeToVehicle" dur="2s" repeatCount="indefinite" />
				</circle>
			</g>
		</svg>
	</div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from "vue";
import formatter, { POWER_UNIT } from "@/mixins/formatter";
import BatteryIcon from "./BatteryIcon.vue";
import VehicleIcon from "../VehicleIcon";
import "@h2d2/shopicons/es/regular/sun";
import "@h2d2/shopicons/es/regular/home";
import "@h2d2/shopicons/es/regular/powersupply";
import type { UiLoadpoint, Battery } from "@/types/evcc";

export default defineComponent({
	name: "VisualEnergyflow",
	components: { BatteryIcon, VehicleIcon },
	mixins: [formatter],
	props: {
		gridPower: { type: Number, default: 0 },
		homePower: { type: Number, default: 0 },
		pvPower: { type: Number, default: 0 },
		batteryConfigured: { type: Boolean, default: false },
		battery: { type: Object as PropType<Battery> },
		loadpoints: { type: Array as PropType<UiLoadpoint[]>, default: () => [] },
	},
	data() {
		return {
			paths: {
				gridToHome: "",
				homeToGrid: "",
				pvToHome: "",
				homeToBattery: "",
				batteryToHome: "",
				homeToVehicle: "",
			},
			resizeObserver: null as ResizeObserver | null,
		};
	},
	computed: {
		powerUnit() {
			const watt = Math.max(this.gridImport, this.pvProduction, this.batteryDischarge, this.gridExport);
			if (watt >= 1_000_000) return POWER_UNIT.MW;
			if (watt >= 1000) return POWER_UNIT.KW;
			return POWER_UNIT.W;
		},
		gridImport() {
			return Math.max(0, this.gridPower);
		},
		gridExport() {
			return Math.max(0, this.gridPower * -1);
		},
		pvProduction() {
			return Math.abs(this.pvPower);
		},
		batteryPower() {
			return this.battery?.power ?? 0;
		},
		batterySoc() {
			return this.battery?.soc ?? 0;
		},
		batteryCharge() {
			return this.batteryPower < 0 ? Math.abs(this.batteryPower) : 0;
		},
		batteryDischarge() {
			return this.batteryPower > 0 ? this.batteryPower : 0;
		},
		loadpointsPower() {
			return this.loadpoints.reduce((sum, lp) => sum + (lp.chargePower || 0), 0);
		},
		activeLoadpointsCount() {
			return this.loadpoints.filter((lp) => lp.charging).length;
		},
		vehicleIcons() {
			if (this.activeLoadpointsCount > 0) {
				return this.loadpoints.filter(lp => lp.charging).map((lp) => lp.icon);
			}
			return ["car"];
		},
	},
	mounted() {
		// Wait for DOM to be fully rendered
		this.$nextTick(() => {
			// Small delay to ensure CSS layout is fully settled
			setTimeout(() => {
				this.calculatePaths();
			}, 100);

			// Re-calculate on window resize or dynamically to keep paths aligned with nodes
			this.resizeObserver = new ResizeObserver(() => {
				this.calculatePaths();
			});
			// Also watch for updates to dynamically recalculate paths in case of reflows
			this.$watch(() => [this.gridImport, this.gridExport, this.pvProduction, this.homePower, this.batteryPower, this.loadpointsPower], () => {
				this.calculatePaths();
			}, { deep: true });
			if (this.$refs.container) {
				this.resizeObserver.observe(this.$refs.container as Element);
			}
		});
	},
	beforeUnmount() {
		if (this.resizeObserver && this.$refs.container) {
			this.resizeObserver.unobserve(this.$refs.container as Element);
		}
	},
	methods: {
		fmtPercentage(val: number) {
			return val ? Math.round(val) + "%" : "0%";
		},
		calculatePaths() {
			// Get actual DOM elements
			const container = this.$refs.container as HTMLElement;
			let home = this.$refs.homeNode;
			if (Array.isArray(home)) home = home[0];
			// Unpack from proxy if necessary (vue3 reactive refs)
			if (home && "$el" in home) home = home.$el;

			if (!container || !home || !(home instanceof HTMLElement)) return;

			const containerRect = container.getBoundingClientRect();

			const getCenter = (el: HTMLElement) => {
				const rect = el.getBoundingClientRect();
				// Calculate relative to the container
				return {
					x: rect.left - containerRect.left + rect.width / 2,
					y: rect.top - containerRect.top + rect.height / 2,
				};
			};

			const cHome = getCenter(home);
			// Offset radius from center so lines don't draw over the circle
			const r = 55;

			// Grid
			let grid = this.$refs.gridNode;
			if (Array.isArray(grid)) grid = grid[0];
			if (grid && "$el" in grid) grid = grid.$el;

			if (grid && grid instanceof HTMLElement) {
				const cGrid = getCenter(grid);
				// Calculate points on the circumference
				const angle = Math.atan2(cHome.y - cGrid.y, cHome.x - cGrid.x);
				const p1 = { x: cGrid.x + Math.cos(angle) * r, y: cGrid.y + Math.sin(angle) * r };
				const p2 = { x: cHome.x - Math.cos(angle) * r, y: cHome.y - Math.sin(angle) * r };

				this.paths.gridToHome = `M ${p1.x} ${p1.y} L ${p2.x} ${p2.y}`;
				this.paths.homeToGrid = `M ${p2.x} ${p2.y} L ${p1.x} ${p1.y}`;
			}

			// PV
			let pv = this.$refs.pvNode;
			if (Array.isArray(pv)) pv = pv[0];
			if (pv && "$el" in pv) pv = pv.$el;

			if (pv && pv instanceof HTMLElement) {
				const cPv = getCenter(pv);
				const angle = Math.atan2(cHome.y - cPv.y, cHome.x - cPv.x);
				const p1 = { x: cPv.x + Math.cos(angle) * r, y: cPv.y + Math.sin(angle) * r };
				const p2 = { x: cHome.x - Math.cos(angle) * r, y: cHome.y - Math.sin(angle) * r };

				this.paths.pvToHome = `M ${p1.x} ${p1.y} L ${p2.x} ${p2.y}`;
			}

			// Battery
			let battery = this.$refs.batteryNode;
			if (Array.isArray(battery)) battery = battery[0];
			if (battery && "$el" in battery) battery = battery.$el;

			if (battery && battery instanceof HTMLElement) {
				const cBattery = getCenter(battery);
				const angle = Math.atan2(cBattery.y - cHome.y, cBattery.x - cHome.x);
				const p1 = { x: cHome.x + Math.cos(angle) * r, y: cHome.y + Math.sin(angle) * r };
				const p2 = { x: cBattery.x - Math.cos(angle) * r, y: cBattery.y - Math.sin(angle) * r };

				this.paths.homeToBattery = `M ${p1.x} ${p1.y} L ${p2.x} ${p2.y}`;
				this.paths.batteryToHome = `M ${p2.x} ${p2.y} L ${p1.x} ${p1.y}`;
			}

			// Vehicle
			let vehicle = this.$refs.vehicleNode;
			if (Array.isArray(vehicle)) vehicle = vehicle[0];
			if (vehicle && "$el" in vehicle) vehicle = vehicle.$el;

			if (vehicle && vehicle instanceof HTMLElement) {
				const cVehicle = getCenter(vehicle);
				const angle = Math.atan2(cVehicle.y - cHome.y, cVehicle.x - cHome.x);
				const p1 = { x: cHome.x + Math.cos(angle) * r, y: cHome.y + Math.sin(angle) * r };
				const p2 = { x: cVehicle.x - Math.cos(angle) * r, y: cVehicle.y - Math.sin(angle) * r };

				this.paths.homeToVehicle = `M ${p1.x} ${p1.y} L ${p2.x} ${p2.y}`;
			}
		}
	}
});
</script>

<style scoped>
.visual-energyflow {
	position: relative;
	width: 100%;
	min-height: 450px;
	padding: 3rem 0;
	background: var(--evcc-box);
	border-radius: 1rem;
	overflow: hidden;
}

.node-container {
	display: flex;
	flex-direction: column;
	justify-content: space-between;
	height: 100%;
	position: relative;
	z-index: 2;
}

.row {
	display: flex;
	justify-content: space-around;
	width: 100%;
	margin-bottom: 3rem;
}

.center-row {
	justify-content: center;
}

.node {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	width: 120px;
	height: 120px;
	border-radius: 50%;
	background-color: var(--evcc-background);
	box-shadow: 0 4px 6px rgba(0,0,0,0.1);
	border: 3px solid var(--evcc-gray);
	transition: all 0.3s ease;
	position: relative;
	text-align: center;
	padding: 0.5rem;
	box-sizing: border-box;
}

.grid-node.active { border-color: var(--evcc-grid); }
.pv-node.active { border-color: var(--evcc-pv); }
.home-node.active { border-color: var(--evcc-gray); }
.battery-node.active { border-color: var(--evcc-battery); }
.vehicle-node.active { border-color: var(--evcc-export); }

.icon {
	font-size: 2rem;
	margin-bottom: 0.15rem;
}

.label {
	font-size: 0.8rem;
	font-weight: bold;
	color: var(--evcc-default-text);
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	width: 100%;
}

.value {
	font-size: 0.75rem;
	color: var(--evcc-gray);
	text-align: center;
	word-wrap: break-word;
	line-height: 1.1;
	margin-top: 0.15rem;
}

.flow-lines {
	position: absolute;
	top: 0;
	left: 0;
	z-index: 1;
	pointer-events: none;
}

.flow-path {
	fill: none;
	stroke: var(--evcc-gray);
	stroke-width: 4;
	stroke-dasharray: 8 8;
	opacity: 0.4;
	transition: all 0.3s ease;
}

.flow-path.grid-in { stroke: var(--evcc-grid); opacity: 0.9; stroke-dasharray: none; stroke-width: 5; }
.flow-path.grid-out { stroke: var(--evcc-export); opacity: 0.9; stroke-dasharray: none; stroke-width: 5; }
.flow-path.pv-in { stroke: var(--evcc-pv); opacity: 0.9; stroke-dasharray: none; stroke-width: 5; }
.flow-path.battery-charge { stroke: var(--evcc-battery); opacity: 0.9; stroke-dasharray: none; stroke-width: 5; }
.flow-path.battery-discharge { stroke: var(--evcc-battery); opacity: 0.9; stroke-dasharray: none; stroke-width: 5; }
.flow-path.vehicle-in { stroke: var(--evcc-export); opacity: 0.9; stroke-dasharray: none; stroke-width: 5; }

.dot {
	fill: var(--evcc-gray);
}
.grid-dot-in { fill: var(--evcc-grid); }
.grid-dot-out { fill: var(--evcc-export); }
.pv-dot { fill: var(--evcc-pv); }
.battery-dot { fill: var(--evcc-battery); }
.vehicle-dot { fill: var(--evcc-export); }
</style>
