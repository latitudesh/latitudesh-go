package latitude

import (
	"testing"
)

const (
	testFirewallType = "firewalls"
)

func deleteFirewall(t *testing.T, c *Client, id string) {
	if _, err := c.Firewalls.Delete(id); err != nil {
		t.Fatal(err)
	}
}

func TestAccFirewallBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, stopRecord := setup(t)
	defer stopRecord()

	// Setup a project for testing
	projectName := testProjectPrefix + randString8()
	pcr := ProjectCreateRequest{
		Data: ProjectCreateData{
			Type: "projects",
			Attributes: ProjectCreateAttributes{
				Name:             projectName,
				Environment:      "Development",
				ProvisioningType: "reserved",
			},
		},
	}
	project, _, err := c.Projects.Create(&pcr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteProject(t, c, project.ID)

	var firewallID string
	t.Run("Create Firewall", func(t *testing.T) {
		// Create a new firewall
		firewallName := "test-firewall-" + randString8()
		fcr := FirewallCreateRequest{
			Data: FirewallCreateData{
				Type: testFirewallType,
				Attributes: FirewallCreateAttributes{
					Name:    firewallName,
					Project: project.Slug,
					Rules: []FirewallRule{
						{
							From:     "192.168.1.1",
							To:       "192.168.1.10",
							Port:     "80",
							Protocol: "TCP",
						},
						{
							From:     "192.168.1.5",
							To:       "192.168.1.15",
							Port:     "443",
							Protocol: "TCP",
						},
					},
				},
			},
		}
		fw, _, err := c.Firewalls.Create(&fcr)
		if err != nil {
			t.Fatal(err)
		}

		firewallID = fw.ID

		assertEqual(t, fw.Name, firewallName, "Firewall Name")
		assertEqual(t, len(fw.Rules), 3, "Firewall Rules") // 2 custom + 1 default SSH rule
	})

	defer deleteFirewall(t, c, firewallID)

	var firewallName string
	t.Run("Update Firewall", func(t *testing.T) {
		newName := "updated-firewall-" + randString8()
		fur := FirewallUpdateRequest{
			Data: FirewallUpdateData{
				ID:   firewallID,
				Type: testFirewallType,
				Attributes: FirewallUpdateAttributes{
					Name: newName,
					Rules: []FirewallRule{
						{
							From:     "192.168.1.1",
							To:       "192.168.1.10",
							Port:     "80",
							Protocol: "TCP",
						},
						{
							From:     "10.0.0.1",
							To:       "10.0.0.10",
							Port:     "8080",
							Protocol: "TCP",
						},
					},
				},
			},
		}

		fw, _, err := c.Firewalls.Update(firewallID, &fur)
		if err != nil {
			t.Fatal(err)
		}

		firewallName = fw.Name
		assertEqual(t, firewallName, newName, "Updated Firewall Name")
		assertEqual(t, len(fw.Rules), 2, "Updated Firewall Rules")
	})

	t.Run("Get Firewall", func(t *testing.T) {
		gotFirewall, _, err := c.Firewalls.Get(firewallID, nil)
		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, gotFirewall.Name, firewallName, "Firewall Name")
	})

	t.Run("List Firewalls", func(t *testing.T) {
		firewalls, _, err := c.Firewalls.List(nil)
		if err != nil {
			t.Fatal(err)
		}

		if len(firewalls) == 0 {
			t.Fatalf("Firewall List should contain at least one firewall")
		}
	})

	// Test server assignment - we need a server first
	var serverID string
	t.Run("Create Server for Assignment", func(t *testing.T) {
		// First get regions, plans, and operating systems for server creation
		regions, _, err := c.Regions.List(nil)
		if err != nil || len(regions) == 0 {
			t.Skip("Skipping assignment test - no regions available")
		}

		plans, _, err := c.Plans.List(nil)
		if err != nil || len(plans) == 0 {
			t.Skip("Skipping assignment test - no plans available")
		}

		oses, _, err := c.OperatingSystems.List(nil)
		if err != nil || len(oses) == 0 {
			t.Skip("Skipping assignment test - no operating systems available")
		}

		scr := ServerCreateRequest{
			Data: ServerCreateData{
				Type: "servers",
				Attributes: ServerCreateAttributes{
					Hostname:        "fw-test-server-" + randString8(),
					Plan:            plans[0].Slug,
					Site:            regions[0].Slug,
					OperatingSystem: oses[0].Slug,
					Project:         project.Slug,
				},
			},
		}

		server, _, err := c.Servers.Create(&scr)
		if err != nil {
			t.Skipf("Skipping assignment test - could not create server: %v", err)
		}

		serverID = server.ID
	})

	if serverID != "" {
		defer func() {
			// Cleanup the server
			if _, err := c.Servers.Delete(serverID); err != nil {
				t.Logf("Failed to delete test server: %v", err)
			}
		}()

		var assignmentID string
		t.Run("Create Firewall Assignment", func(t *testing.T) {
			acr := FirewallAssignmentCreateRequest{
				Data: FirewallAssignmentCreateData{
					Type: "firewall_server",
					Attributes: FirewallAssignmentCreateAttributes{
						// This will be sent as "server_id" in the JSON due to the struct tag
						Server: serverID,
					},
				},
			}

			assignment, _, err := c.Firewalls.CreateAssignment(firewallID, &acr)
			if err != nil {
				t.Fatal(err)
			}

			assignmentID = assignment.ID
			assertEqual(t, assignment.Server.ID, serverID, "Assignment Server ID")
		})

		t.Run("List Firewall Assignments", func(t *testing.T) {
			assignments, _, err := c.Firewalls.ListAssignments(firewallID, nil)
			if err != nil {
				t.Fatal(err)
			}

			if len(assignments) == 0 {
				t.Fatalf("Firewall Assignments List should contain at least one assignment")
			}
		})

		t.Run("Delete Firewall Assignment", func(t *testing.T) {
			_, err := c.Firewalls.DeleteAssignment(firewallID, assignmentID)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
