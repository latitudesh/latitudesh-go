package latitude

import (
	"testing"
)

func TestAccRegionBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, stopRecord := setup(t)
	defer stopRecord()

	rl, _, err := c.Regions.List(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(rl) == 0 {
		t.Fatalf("Region List should contain at least one plan")
	}

	// Get first listed region
	gotRegion, _, err := c.Regions.Get(rl[0].ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Check region data
	assertEqual(t, gotRegion.ID, rl[0].ID, "Regions ID")
	assertEqual(t, gotRegion.Type, rl[0].Type, "Regions Type")
	assertEqual(t, gotRegion.Slug, rl[0].Slug, "Regions Slug")
	assertEqual(t, gotRegion.Facility, rl[0].Facility, "Regions Facility")
	assertEqual(t, gotRegion.CountryName, rl[0].CountryName, "Regions Country Name")
	assertEqual(t, gotRegion.CountrySlug, rl[0].CountrySlug, "Regions Country Slug")
}

func TestAccRegionFilter(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, stopRecord := setup(t)
	defer stopRecord()

	rl, _, err := c.Regions.List(new(GetOptions).Filter("slug", testRegionDefault))
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, len(rl), 1, "Region List")
}
