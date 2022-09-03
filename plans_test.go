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
	if gotPlan.ID != pl[0].ID {
		t.Fatalf("Expected the id of the GOT plan to be %s, not %s", pl[0].ID, gotPlan.ID)
	}
	if gotPlan.Line != pl[0].Line {
		t.Fatalf("Expected the line of the GOT plan to be %s, not %s", pl[0].Line, gotPlan.Line)
	}
	if gotPlan.Name != pl[0].Name {
		t.Fatalf("Expected the name of the GOT plan to be %s, not %s", pl[0].Name, gotPlan.Name)
	}
	if gotPlan.Type != pl[0].Type {
		t.Fatalf("Expected the type of the GOT plan to be %s, not %s", pl[0].Type, gotPlan.Type)
	}

	// Check plan features
	if gotPlan.Features.RAID != pl[0].Features.RAID {
		t.Fatalf("Expected the RAID of the GOT plan to be %t, not %t", pl[0].Features.RAID, gotPlan.Features.RAID)
	}
	if gotPlan.Features.SSH != pl[0].Features.SSH {
		t.Fatalf("Expected the SSH of the GOT plan to be %t, not %t", pl[0].Features.SSH, gotPlan.Features.SSH)
	}
	if gotPlan.Features.UserData != pl[0].Features.UserData {
		t.Fatalf("Expected the user data of the GOT plan to be %t, not %t", pl[0].Features.UserData, gotPlan.Features.UserData)
	}
	if gotPlan.Specs.Memory.Total != pl[0].Specs.Memory.Total {
		t.Fatalf("Expected the memory total of the GOT plan to be %v, not %v", pl[0].Specs.Memory.Total, gotPlan.Specs.Memory.Total)
	}

	// Check plan specs (CPUs)
	for i, cpu := range gotPlan.Specs.CPUs {
		if cpu.Type != pl[0].Specs.CPUs[i].Type {
			t.Fatalf("Expected the CPU type of the GOT plan to be %s, not %s", cpu.Type, pl[0].Specs.CPUs[i].Type)
		}
		if cpu.Clock != pl[0].Specs.CPUs[i].Clock {
			t.Fatalf("Expected the CPU clock of the GOT plan to be %f, not %f", cpu.Clock, pl[0].Specs.CPUs[i].Clock)
		}
		if cpu.Cores != pl[0].Specs.CPUs[i].Cores {
			t.Fatalf("Expected the CPU cores of the GOT plan to be %d, not %d", cpu.Cores, pl[0].Specs.CPUs[i].Cores)
		}
		if cpu.Count != pl[0].Specs.CPUs[i].Count {
			t.Fatalf("Expected the CPU count of the GOT plan to be %d, not %d", cpu.Count, pl[0].Specs.CPUs[i].Count)
		}
	}

	// Check plan specs (drives)
	for i, drive := range gotPlan.Specs.Drives {
		if drive.Type != pl[0].Specs.Drives[i].Type {
			t.Fatalf("Expected the drive type of the GOT plan to be %s, not %s", drive.Type, pl[0].Specs.Drives[i].Type)
		}
		if drive.Count != pl[0].Specs.Drives[i].Count {
			t.Fatalf("Expected the drive count of the GOT plan to be %d, not %d", drive.Count, pl[0].Specs.Drives[i].Count)
		}
		if drive.Size != pl[0].Specs.Drives[i].Size {
			t.Fatalf("Expected the drive size of the GOT plan to be %s, not %s", drive.Size, pl[0].Specs.Drives[i].Size)
		}
	}

	// Check plan specs (NICs)
	for i, nic := range gotPlan.Specs.NICs {
		if nic.Count != pl[0].Specs.NICs[i].Count {
			t.Fatalf("Expected the NIC count of the GOT plan to be %d, not %d", nic.Count, pl[0].Specs.NICs[i].Count)
		}
		if nic.Type != pl[0].Specs.NICs[i].Type {
			t.Fatalf("Expected the NIC type of the GOT plan to be %s, not %s", nic.Type, pl[0].Specs.NICs[i].Type)
		}
	}

	// Check plan availability
	for i, availability := range gotPlan.Availibility {
		if availability.Region.ID != pl[0].Availibility[i].Region.ID {
			t.Fatalf("Expected the availability region id of the GOT plan to be %d, not %d", availability.Region.ID, pl[0].Availibility[i].Region.ID)
		}
		if availability.Region.Name != pl[0].Availibility[i].Region.Name {
			t.Fatalf("Expected the availability region name of the GOT plan to be %s, not %s", availability.Region.Name, pl[0].Availibility[i].Region.Name)
		}
	}
}
