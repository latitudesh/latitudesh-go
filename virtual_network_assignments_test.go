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

			assertEqual(t, va.Type, vaTest.Type, "Vlan Assignment Type")
			assertEqual(t, va.Vid, vaTest.Vid, "Vlan Assignment Vid")
			assertEqual(t, va.Description, vaTest.Description, "Vlan Assignment Description")
			assertEqual(t, va.VirtualNetworkID, vaTest.VirtualNetworkID, "Vlan Assignment VirtualNetworkID")
			assertEqual(t, va.Status, vaTest.Status, "Vlan Assignment Status")
			assertEqual(t, va.ServerID, vaTest.ServerID, "Vlan Assignment ServerID")
			assertEqual(t, va.ServerHostname, vaTest.ServerHostname, "Vlan Assignment ServerHostname")
			assertEqual(t, va.ServerLabel, vaTest.ServerLabel, "Vlan Assignment ServerLabel")
			assertEqual(t, va.ServerStatus, vaTest.ServerStatus, "Vlan Assignment ServerStatus")

			return
		}
		t.Fatalf("Vlan Assignment with id %s not found", vaTest.ID)
	})
}
