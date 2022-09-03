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
	if gotPlan.Data.ID != pl[0].ID {
		t.Fatalf("Expected the id of the GOT plan to be %s, not %s", pl[0].ID, gotPlan.Data.ID)
	}
}
