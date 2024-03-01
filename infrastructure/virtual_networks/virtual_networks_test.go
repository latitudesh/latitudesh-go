package virtual_networks_test

import (
	"testing"

	latitude "github.com/latitudesh/latitudesh-go"
	vnet "github.com/latitudesh/latitudesh-go/infrastructure/virtual_networks"
)

func TestAccVirtualNetworkBasic(t *testing.T) {
	latitude.SkipUnlessAcceptanceTestsAllowed(t)
	c, projectID, teardown := latitude.SetupWithProject(t)
	defer teardown()

	// Create Virtual Network
	createRequest := vnet.VirtualNetworkCreateRequest{
		Data: vnet.VirtualNetworkCreateData{
			Type: "virtual_network",
			Attributes: vnet.VirtualNetworkCreateAttributes{
				Description: "Testing Virtual Network via golang client",
				Site:        latitude.TestSite(),
				Project:     projectID,
			},
		},
	}

	vnNew, _, err := c.VirtualNetworks.Create(&createRequest)
	if err != nil {
		t.Fatal(err)
	}
	defer c.VirtualNetworks.Delete(vnNew.ID)

	updateRequest := vnet.VirtualNetworkUpdateRequest{
		Data: vnet.VirtualNetworkUpdateData{
			ID:   vnNew.ID,
			Type: "virtual_networks",
			Attributes: vnet.VirtualNetworkUpdateAttributes{
				Description: "Updating Virtual Network via golang client",
			},
		},
	}

	_, _, err = c.VirtualNetworks.Update(vnNew.ID, &updateRequest)
	if err != nil {
		t.Fatal(err)
	}

	vnTest, _, err := c.VirtualNetworks.Get(vnNew.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	vnl, _, err := c.VirtualNetworks.List(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(vnl) == 0 {
		t.Fatalf("Virtual Network List should contain at least one virtual network")
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
			t.Fatalf("Expected the site id of the Virtual Network to be %s, not %s", vnTest.SiteId, vn.SiteId)
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
		return
	}
	t.Fatalf("Virtual Network with id %s not found", vnTest.ID)
}
