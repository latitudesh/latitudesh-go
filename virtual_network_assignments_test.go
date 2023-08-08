package latitude

import (
	"testing"
)

func TestAccVlanAssignmentBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, stopRecord := setup(t)
	defer stopRecord()

	val, _, err := c.VlanAssignments.List(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(val) == 0 {
		t.Fatalf("Vlan Assignment List should contain at least one virtual network")
	}

	vaTest := VlanAssignment{
		ID:               "1189",
		Type:             "virtual_network_assignment",
		VirtualNetworkID: 2054,
		Vid:              2011,
		Description:      "ceph",
		Status:           "connected",
		ServerID:         22892,
		ServerHostname:   "leo-ceph-2",
		ServerLabel:      "224S602585",
		ServerStatus:     "on",
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
			t.Fatalf("Expected the virtual network id of the Vlan Assignment to be %d, not %d", vaTest.VirtualNetworkID, va.VirtualNetworkID)
		}
		if vaTest.Status != va.Status {
			t.Fatalf("Expected the status of the Vlan Assignment to be %s, not %s", vaTest.Status, va.Status)
		}
		if vaTest.ServerID != va.ServerID {
			t.Fatalf("Expected the server id of the Vlan Assignment to be %d, not %d", vaTest.ServerID, va.ServerID)
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
