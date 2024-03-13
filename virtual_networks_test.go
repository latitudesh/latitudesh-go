package latitude

import (
	"testing"
)

func deleteVirtualNetwork(t *testing.T, c *Client, id string) {
	if _, err := c.VirtualNetworks.Delete(id); err != nil {
		t.Fatal(err)
	}
}

func TestAccVirtualNetworkBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	tagIDs, deleteTags := setupTestTags(t, c)
	defer deleteTags()

	var vnID string

	t.Run("Create Virtual Network", func(t *testing.T) {
		createRequest := VirtualNetworkCreateRequest{
			Data: VirtualNetworkCreateData{
				Type: "virtual_network",
				Attributes: VirtualNetworkCreateAttributes{
					Description: "Testing Virtual Network via golang client",
					Site:        testSite(),
					Project:     projectID,
				},
			},
		}

		vnNew, _, err := c.VirtualNetworks.Create(&createRequest)
		if err != nil {
			t.Fatal(err)
		}
		vnID = vnNew.ID
	})
	defer deleteVirtualNetwork(t, c, vnID)

	t.Run("Update Virtual Network", func(t *testing.T) {
		updateRequest := VirtualNetworkUpdateRequest{
			Data: VirtualNetworkUpdateData{
				ID:   vnID,
				Type: "virtual_networks",
				Attributes: VirtualNetworkUpdateAttributes{
					Description: "Updating Virtual Network via golang client",
					Tags:        tagIDs,
				},
			},
		}

		_, _, err := c.VirtualNetworks.Update(vnID, &updateRequest)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Get and List Virtual Networks", func(t *testing.T) {
		vnTest, _, err := c.VirtualNetworks.Get(vnID, nil)
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

			assertEqual(t, vn.Type, vnTest.Type, "Virtual Network Type")
			assertEqual(t, vn.Vid, vnTest.Vid, "Virtual Network Vid")
			assertEqual(t, vn.Description, vnTest.Description, "Virtual Network Description")
			assertEqual(t, vn.City, vnTest.City, "Virtual Network City")
			assertEqual(t, vn.Country, vnTest.Country, "Virtual Network Country")
			assertEqual(t, vn.SiteId, vnTest.SiteId, "Virtual Network SiteId")
			assertEqual(t, vn.SiteName, vnTest.SiteName, "Virtual Network SiteName")
			assertEqual(t, vn.SiteSlug, vnTest.SiteSlug, "Virtual Network SiteSlug")
			assertEqual(t, vn.Facility, vnTest.Facility, "Virtual Network Facility")
			assertEqual(t, len(vn.Tags), 2, "Virtual Network Tags")

			return
		}
		t.Fatalf("Virtual Network with id %s not found", vnTest.ID)
	})
}
