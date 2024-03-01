package regions_test

import (
	"testing"

	latitude "github.com/latitudesh/latitudesh-go"
	api "github.com/latitudesh/latitudesh-go/api_utils"
)

func TestAccRegionBasic(t *testing.T) {
	latitude.SkipUnlessAcceptanceTestsAllowed(t)
	c, stopRecord := latitude.Setup(t)
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
	if gotRegion.ID != rl[0].ID {
		t.Fatalf("Expected the id of the GOT region to be %s, not %s", rl[0].ID, gotRegion.ID)
	}
	if gotRegion.Type != rl[0].Type {
		t.Fatalf("Expected the type of the GOT region to be %s, not %s", rl[0].Type, gotRegion.Type)
	}
	if gotRegion.Slug != rl[0].Slug {
		t.Fatalf("Expected the slug of the GOT region to be %s, not %s", rl[0].Slug, gotRegion.Slug)
	}
	if gotRegion.Facility != rl[0].Facility {
		t.Fatalf("Expected the line of the GOT region to be %s, not %s", rl[0].Facility, gotRegion.Facility)
	}
	if gotRegion.CountryName != rl[0].CountryName {
		t.Fatalf("Expected the country name of the GOT region to be %s, not %s", rl[0].CountryName, gotRegion.CountryName)
	}
	if gotRegion.CountrySlug != rl[0].CountrySlug {
		t.Fatalf("Expected the country slug of the GOT region to be %s, not %s", rl[0].CountrySlug, gotRegion.CountrySlug)
	}
}

func TestAccRegionFilter(t *testing.T) {
	latitude.SkipUnlessAcceptanceTestsAllowed(t)
	c, stopRecord := latitude.Setup(t)
	defer stopRecord()

	rl, _, err := c.Regions.List(new(api.GetOptions).Filter("slug", latitude.TestRegionDefault))
	if err != nil {
		t.Fatal(err)
	}

	if len(rl) != 1 {
		t.Fatalf("Filtered region list should contain one plan: returned %d", len(rl))
	}
}
