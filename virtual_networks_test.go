package latitude

import (
	"testing"
)

func TestAccVirtualNetworkBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, stopRecord := setup(t)
	defer stopRecord()

	vnl, _, err := c.VirtualNetworks.List(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(vnl) == 0 {
		t.Fatalf("Virtual Network List should contain at least one virtual network")
	}

	vnTest := VirtualNetwork{
		ID:               "2054",
		Type:             "virtual_networks",
		Vid:              2011,
		Description:      "ceph",
		City:             "Dallas",
		Country:          "United States",
		SiteId:           3,
		SiteName:         "Dallas",
		SiteSlug:         "DAL",
		Facility:         "Cologix",
		AssignmentsCount: 0,
	}

	// Check Virtual Network data
	for _, vn := range vnl {
		if vn.ID != vnTest.ID {
			continue
		}

		if vnTest.Type != vn.Type {
			t.Fatalf("Expected the type of the Virtual Network to be %s, not %s", vnTest.Type, vn.Type)
		}
		if vnTest.Vid != vn.Vid {
			t.Fatalf("Expected the vid of the Virtual Network to be %d, not %d", vnTest.Vid, vn.Vid)
		}
		if vnTest.Description != vn.Description {
			t.Fatalf("Expected the description of the Virtual Network to be %s, not %s", vnTest.Description, vn.Description)
		}
		if vnTest.City != vn.City {
			t.Fatalf("Expected the region city of the Virtual Network to be %s, not %s", vnTest.City, vn.City)
		}
		if vnTest.Country != vn.Country {
			t.Fatalf("Expected the region country of the Virtual Network to be %s, not %s", vnTest.Country, vn.Country)
		}
		if vnTest.SiteId != vn.SiteId {
			t.Fatalf("Expected the site id of the Virtual Network to be %d, not %d", vnTest.SiteId, vn.SiteId)
		}
		if vnTest.SiteName != vn.SiteName {
			t.Fatalf("Expected the site name of the Virtual Network to be %s, not %s", vnTest.SiteName, vn.SiteName)
		}
		if vnTest.SiteSlug != vn.SiteSlug {
			t.Fatalf("Expected the site slug of the Virtual Network to be %s, not %s", vnTest.SiteSlug, vn.SiteSlug)
		}
		if vnTest.Facility != vn.Facility {
			t.Fatalf("Expected the facility of the Virtual Network to be %s, not %s", vnTest.Facility, vn.Facility)
		}
		if vnTest.AssignmentsCount != vn.AssignmentsCount {
			t.Fatalf("Expected the Assignment count of the Virtual Network to be %d, not %d", vnTest.AssignmentsCount, vn.AssignmentsCount)
		}
		return
	}
	t.Fatalf("Virtual Network with id %s not found", vnTest.ID)
}
