package virtual_network_assignment_test

import (
	"testing"

	latitude "github.com/latitudesh/latitudesh-go"
	servers "github.com/latitudesh/latitudesh-go/infrastructure/servers"
	vnet "github.com/latitudesh/latitudesh-go/infrastructure/virtual_networks"
	vlanassign "github.com/latitudesh/latitudesh-go/infrastructure/virtual_networks_assignments"
)

func deleteServer(t *testing.T, c *latitude.Client, id string) {
	if _, err := c.Servers.Delete(id); err != nil {
		t.Fatal(err)
	}
}

func TestAccVlanAssignmentBasic(t *testing.T) {
	latitude.SkipUnlessAcceptanceTestsAllowed(t)
	c, projectID, teardown := latitude.SetupWithProject(t)
	defer teardown()

	hn := latitude.RandString8()
	scr := servers.ServerCreateRequest{
		Data: servers.ServerCreateData{
			Type: "servers",
			Attributes: servers.ServerCreateAttributes{
				Project:         projectID,
				Plan:            latitude.TestPlan(),
				Site:            latitude.TestSite(),
				OperatingSystem: latitude.TestOperatingSystem(),
				Hostname:        hn,
			},
		},
	}

	s, _, err := c.Servers.Create(&scr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteServer(t, c, s.ID)

	createRequest := vnet.VirtualNetworkCreateRequest{
		Data: vnet.VirtualNetworkCreateData{
			Type: "virtual_network",
			Attributes: vnet.VirtualNetworkCreateAttributes{
				Description: "Testing golang client",
				Site:        latitude.TestSite(),
				Project:     projectID,
			},
		},
	}
	vn, _, err := c.VirtualNetworks.Create(&createRequest)
	if err != nil {
		t.Fatal(err)
	}
	defer c.VirtualNetworks.Delete(vn.ID)

	assignRequest := vlanassign.VlanAssignRequest{
		Data: vlanassign.VlanAssignData{
			Type: "virtual_network_assignment",
			Attributes: vlanassign.VlanAssignAttributes{
				ServerID:         s.ID,
				VirtualNetworkID: vn.ID,
			},
		},
	}

	assign, _, err := c.VlanAssignments.Assign(&assignRequest)
	if err != nil {
		t.Fatal(err)
	}
	defer c.VlanAssignments.Delete(assign.ID)

	vaTest, _, err := c.VlanAssignments.Get(assign.ID)
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
}
