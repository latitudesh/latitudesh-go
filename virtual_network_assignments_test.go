package latitude

import (
	"testing"
)

func deleteVlanAssignment(t *testing.T, c *Client, id string) {
	if _, err := c.VlanAssignments.Delete(id); err != nil {
		t.Fatal(err)
	}
}

func createServer(t *testing.T, c *Client, projectID string) *Server {
	hn := randString8()
	scr := ServerCreateRequest{
		Data: ServerCreateData{
			Type: testServerType,
			Attributes: ServerCreateAttributes{
				Project:         projectID,
				Plan:            testPlan(),
				Site:            testSite(),
				OperatingSystem: testOperatingSystem(),
				Hostname:        hn,
			},
		},
	}

	s, _, err := c.Servers.Create(&scr)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func createVirtualNetwork(t *testing.T, c *Client, projectID string) *VirtualNetwork {
	createRequest := VirtualNetworkCreateRequest{
		Data: VirtualNetworkCreateData{
			Type: "virtual_network",
			Attributes: VirtualNetworkCreateAttributes{
				Description: "Testing golang client",
				Site:        testSite(),
				Project:     projectID,
			},
		},
	}
	vn, _, err := c.VirtualNetworks.Create(&createRequest)
	if err != nil {
		t.Fatal(err)
	}
	return vn
}

func TestAccVlanAssignmentBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	s := createServer(t, c, projectID)
	defer deleteServer(t, c, s.ID)

	vn := createVirtualNetwork(t, c, projectID)
	defer c.VirtualNetworks.Delete(vn.ID)

	var vlanID string

	t.Run("Assing Virtual Network", func(t *testing.T) {
		assignRequest := VlanAssignRequest{
			Data: VlanAssignData{
				Type: "virtual_network_assignment",
				Attributes: VlanAssignAttributes{
					ServerID:         s.ID,
					VirtualNetworkID: vn.ID,
				},
			},
		}

		assign, _, err := c.VlanAssignments.Assign(&assignRequest)
		if err != nil {
			t.Fatal(err)
		}
		vlanID = assign.ID
	})
	defer deleteVlanAssignment(t, c, vlanID)

	t.Run("Get and List Assignments", func(t *testing.T) {
		vaTest, _, err := c.VlanAssignments.Get(vlanID)
		if err != nil {
			t.Fatal(err)
		}

		val, _, err := c.VlanAssignments.List(nil)
		if err != nil {
			t.Fatal(err)
		}

		if len(val) == 0 {
			t.Fatalf("Vlan Assignment List should contain at least one virtual network")
		}

		// Check Vlan Assignment data
		for _, va := range val {
			if va.ID != vaTest.ID {
				continue
			}

			if vaTest.Type != va.Type {
				t.Fatalf("Expected the type of the Vlan Assignment to be %s, not %s", vaTest.Type, va.Type)
			}
			if vaTest.Vid != va.Vid {
				t.Fatalf("Expected the vid of the Vlan Assignment to be %d, not %d", vaTest.Vid, va.Vid)
			}
			if vaTest.Description != va.Description {
				t.Fatalf("Expected the description of the Vlan Assignment to be %s, not %s", vaTest.Description, va.Description)
			}
			if vaTest.VirtualNetworkID != va.VirtualNetworkID {
				t.Fatalf("Expected the virtual network id of the Vlan Assignment to be %s, not %s", vaTest.VirtualNetworkID, va.VirtualNetworkID)
			}
			if vaTest.Status != va.Status {
				t.Fatalf("Expected the status of the Vlan Assignment to be %s, not %s", vaTest.Status, va.Status)
			}
			if vaTest.ServerID != va.ServerID {
				t.Fatalf("Expected the server id of the Vlan Assignment to be %s, not %s", vaTest.ServerID, va.ServerID)
			}
			if vaTest.ServerHostname != va.ServerHostname {
				t.Fatalf("Expected the server hostname of the Vlan Assignment to be %s, not %s", vaTest.ServerHostname, va.ServerHostname)
			}
			if vaTest.ServerLabel != va.ServerLabel {
				t.Fatalf("Expected the server label of the Vlan Assignment to be %s, not %s", vaTest.ServerLabel, va.ServerLabel)
			}
			if vaTest.ServerStatus != va.ServerStatus {
				t.Fatalf("Expected the server status of the Vlan Assignment to be %s, not %s", vaTest.ServerStatus, va.ServerStatus)
			}
			return
		}
		t.Fatalf("Vlan Assignment with id %s not found", vaTest.ID)
	})
}
