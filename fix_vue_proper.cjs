const fs = require('fs');
const vuePath = 'assets/js/components/Energyflow/VisualEnergyflow.vue';
let content = fs.readFileSync(vuePath, 'utf8');

content = content.replace(/let home: any = this\.\$refs\.homeNode;/g, "let home: any = (this.$refs as any)['homeNode'];");
content = content.replace(/let grid: any = this\.\$refs\.gridNode;/g, "let grid: any = (this.$refs as any)['gridNode'];");
content = content.replace(/let pv: any = this\.\$refs\.pvNode;/g, "let pv: any = (this.$refs as any)['pvNode'];");
content = content.replace(/let battery: any = this\.\$refs\.batteryNode;/g, "let battery: any = (this.$refs as any)['batteryNode'];");
content = content.replace(/let vehicle: any = this\.\$refs\.vehicleNode;/g, "let vehicle: any = (this.$refs as any)['vehicleNode'];");

content = content.replace(/let home = this\.\$refs\.homeNode;/g, "let home: any = (this.$refs as any)['homeNode'];");
content = content.replace(/let grid = this\.\$refs\.gridNode;/g, "let grid: any = (this.$refs as any)['gridNode'];");
content = content.replace(/let pv = this\.\$refs\.pvNode;/g, "let pv: any = (this.$refs as any)['pvNode'];");
content = content.replace(/let battery = this\.\$refs\.batteryNode;/g, "let battery: any = (this.$refs as any)['batteryNode'];");
content = content.replace(/let vehicle = this\.\$refs\.vehicleNode;/g, "let vehicle: any = (this.$refs as any)['vehicleNode'];");

content = content.replace(/if \(home && "\$el" in home\) home = home\.\$el;/g, "if (home && \"$el\" in home) home = (home as any).$el;");
content = content.replace(/if \(grid && "\$el" in grid\) grid = grid\.\$el;/g, "if (grid && \"$el\" in grid) grid = (grid as any).$el;");
content = content.replace(/if \(pv && "\$el" in pv\) pv = pv\.\$el;/g, "if (pv && \"$el\" in pv) pv = (pv as any).$el;");
content = content.replace(/if \(battery && "\$el" in battery\) battery = battery\.\$el;/g, "if (battery && \"$el\" in battery) battery = (battery as any).$el;");
content = content.replace(/if \(vehicle && "\$el" in vehicle\) vehicle = vehicle\.\$el;/g, "if (vehicle && \"$el\" in vehicle) vehicle = (vehicle as any).$el;");

content = content.replace(/if \(home && "\$el" in \(home as any\)\) home = \(home as any\)\.\$el;/g, "if (home && \"$el\" in home) home = (home as any).$el;");
content = content.replace(/if \(grid && "\$el" in \(grid as any\)\) grid = \(grid as any\)\.\$el;/g, "if (grid && \"$el\" in grid) grid = (grid as any).$el;");
content = content.replace(/if \(pv && "\$el" in \(pv as any\)\) pv = \(pv as any\)\.\$el;/g, "if (pv && \"$el\" in pv) pv = (pv as any).$el;");
content = content.replace(/if \(battery && "\$el" in \(battery as any\)\) battery = \(battery as any\)\.\$el;/g, "if (battery && \"$el\" in battery) battery = (battery as any).$el;");
content = content.replace(/if \(vehicle && "\$el" in \(vehicle as any\)\) vehicle = \(vehicle as any\)\.\$el;/g, "if (vehicle && \"$el\" in vehicle) vehicle = (vehicle as any).$el;");

content = content.replace("const container = this.$refs.container as HTMLElement;", "const container = (this.$refs as any)['container'] as HTMLElement | undefined;\n\t\t\tif (!container) return;");

content = content.replace(/if \(!container \|\| !home \|\| !\(home instanceof HTMLElement\)\) return;/g, "if (!home) return;");

fs.writeFileSync(vuePath, content);
