package latitude

import (
	"testing"
)

func TestAccPlanBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, stopRecord := setup(t)
	defer stopRecord()

	pl, _, err := c.Plans.List(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(pl) == 0 {
		t.Fatalf("Plan List should contain at least one plan")
	}

	// Get first listed plan
	gotPlan, _, err := c.Plans.Get(pl[0].ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Check plan data
	assertEqual(t, gotPlan.ID, pl[0].ID, "Plan ID")
	assertEqual(t, gotPlan.Name, pl[0].Name, "Plan Name")
	assertEqual(t, gotPlan.Type, pl[0].Type, "Plan ID")

	// Check plan features
	assertEqual(t, len(gotPlan.Features), len(pl[0].Features), "Plan features lenght")

	// Check plan specs Memorys
	assertEqual(t, gotPlan.Specs.Memory.Total, pl[0].Specs.Memory.Total, "Plan total memory")

	// Check plan specs (CPUs)
	cpu := gotPlan.Specs.CPU
	assertEqual(t, cpu.Type, pl[0].Specs.CPU.Type, "Plan CPU type")
	assertEqual(t, cpu.Clock, pl[0].Specs.CPU.Clock, "Plan CPU clock")
	assertEqual(t, cpu.Cores, pl[0].Specs.CPU.Cores, "Plan CPU cores")
	assertEqual(t, cpu.Count, pl[0].Specs.CPU.Count, "Plan CPU count")

	// Check plan specs (drives)
	for i, drive := range gotPlan.Specs.Drives {
		assertEqual(t, drive.Type, pl[0].Specs.Drives[i].Type, "Plan Drive type")
		assertEqual(t, drive.Count, pl[0].Specs.Drives[i].Count, "Plan Drive count")
		assertEqual(t, drive.Size, pl[0].Specs.Drives[i].Size, "Plan Drive size")
	}

	// Check plan specs (NICs)
	for i, nic := range gotPlan.Specs.NICs {
		assertEqual(t, nic.Count, pl[0].Specs.NICs[i].Count, "Plan NIC count")
		assertEqual(t, nic.Type, pl[0].Specs.NICs[i].Type, "Plan NIC type")
	}

	// Check plan availability
	for i, stock := range gotPlan.InStock {
		assertEqual(t, stock, pl[0].InStock[i], "Plan available stock")
	}
}

func TestAccPlanFilter(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, stopRecord := setup(t)
	defer stopRecord()

	pl, _, err := c.Plans.List(new(GetOptions).Filter("slug", testPlanDefault))
	if err != nil {
		t.Fatal(err)
	}

	if len(pl) != 1 {
		t.Fatalf("Filtered plan list should contain one plan: returned %d", len(pl))
	}
}
