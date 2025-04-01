package main

import (
	"fmt"
	"log"
	"os"

	latitude "github.com/latitudesh/latitudesh-go"
)

func main() {
	// Set the API token for authentication
	apiToken := os.Getenv("LATITUDE_API_TOKEN")
	if apiToken == "" {
		log.Fatal("LATITUDE_API_TOKEN environment variable must be set")
	}

	// Create a new client
	client := latitude.NewClientWithAuth("Firewall Example", apiToken, nil)

	// Set the project name/slug to use for testing
	projectSlug := os.Getenv("LATITUDE_PROJECT")
	if projectSlug == "" {
		log.Fatal("LATITUDE_PROJECT environment variable must be set")
	}

	// List existing firewalls
	fmt.Println("Listing existing firewalls...")
	firewalls, _, err := client.Firewalls.List(nil)
	if err != nil {
		log.Fatalf("Error listing firewalls: %v", err)
	}

	fmt.Printf("Found %d firewalls\n", len(firewalls))
	for i, fw := range firewalls {
		fmt.Printf("%d. %s (ID: %s) - Project: %s\n", i+1, fw.Name, fw.ID, fw.Project.Name)
	}

	// Create a new firewall
	fmt.Println("\nCreating a new firewall...")
	newFirewall, err := createFirewall(client, projectSlug)
	if err != nil {
		log.Fatalf("Error creating firewall: %v", err)
	}

	fmt.Printf("Created firewall: %s (ID: %s)\n", newFirewall.Name, newFirewall.ID)
	fmt.Println("Rules:")
	for i, rule := range newFirewall.Rules {
		fmt.Printf("  %d. From: %s, To: %s, Port: %s, Protocol: %s\n",
			i+1, rule.From, rule.To, rule.Port, rule.Protocol)
	}

	// Update the firewall
	fmt.Println("\nUpdating the firewall...")
	updatedFirewall, err := updateFirewall(client, newFirewall.ID)
	if err != nil {
		log.Fatalf("Error updating firewall: %v", err)
	}

	fmt.Printf("Updated firewall: %s (ID: %s)\n", updatedFirewall.Name, updatedFirewall.ID)
	fmt.Println("Updated rules:")
	for i, rule := range updatedFirewall.Rules {
		fmt.Printf("  %d. From: %s, To: %s, Port: %s, Protocol: %s\n",
			i+1, rule.From, rule.To, rule.Port, rule.Protocol)
	}

	// List servers to potentially assign to the firewall
	fmt.Println("\nListing servers for potential assignment...")
	servers, _, err := client.Servers.List(projectSlug, nil)
	if err != nil {
		log.Fatalf("Error listing servers: %v", err)
	}

	if len(servers) > 0 {
		fmt.Printf("Found %d servers. First server: %s (ID: %s)\n",
			len(servers), servers[0].Hostname, servers[0].ID)

		// Create firewall assignment if we have servers
		fmt.Println("\nCreating firewall assignment...")
		assignment, err := createFirewallAssignment(client, newFirewall.ID, servers[0].ID)
		if err != nil {
			log.Printf("Error creating firewall assignment: %v", err)
		} else {
			fmt.Printf("Created firewall assignment (ID: %s) for server: %s\n",
				assignment.ID, assignment.Server.Hostname)

			// List assignments
			fmt.Println("\nListing firewall assignments...")
			assignments, _, err := client.Firewalls.ListAssignments(newFirewall.ID, nil)
			if err != nil {
				log.Printf("Error listing firewall assignments: %v", err)
			} else {
				fmt.Printf("Found %d assignments\n", len(assignments))

				// Delete the assignment
				fmt.Println("\nDeleting firewall assignment...")
				_, err = client.Firewalls.DeleteAssignment(newFirewall.ID, assignment.ID)
				if err != nil {
					log.Printf("Error deleting firewall assignment: %v", err)
				} else {
					fmt.Println("Successfully deleted firewall assignment")
				}
			}
		}
	} else {
		fmt.Println("No servers found for assignment testing")
	}

	// Delete the firewall
	fmt.Println("\nDeleting firewall...")
	_, err = client.Firewalls.Delete(newFirewall.ID)
	if err != nil {
		log.Fatalf("Error deleting firewall: %v", err)
	}
	fmt.Println("Successfully deleted firewall")

	fmt.Println("\nFirewall test completed successfully!")
}

// createFirewall creates a new firewall with some example rules
func createFirewall(client *latitude.Client, projectSlug string) (*latitude.Firewall, error) {
	createRequest := &latitude.FirewallCreateRequest{
		Data: latitude.FirewallCreateData{
			Type: "firewalls",
			Attributes: latitude.FirewallCreateAttributes{
				Name:    "example-firewall",
				Project: projectSlug,
				Rules: []latitude.FirewallRule{
					{
						From:     "192.168.1.0/24",
						To:       "ANY",
						Port:     "80",
						Protocol: "TCP",
					},
					{
						From:     "10.0.0.0/8",
						To:       "192.168.1.100",
						Port:     "443",
						Protocol: "TCP",
					},
					{
						From:     "172.16.0.0/12",
						To:       "ANY",
						Port:     "3000-4000",
						Protocol: "UDP",
					},
				},
			},
		},
	}

	firewall, _, err := client.Firewalls.Create(createRequest)
	return firewall, err
}

// updateFirewall updates an existing firewall with new rules
func updateFirewall(client *latitude.Client, firewallID string) (*latitude.Firewall, error) {
	updateRequest := &latitude.FirewallUpdateRequest{
		Data: latitude.FirewallUpdateData{
			ID:   firewallID,
			Type: "firewalls",
			Attributes: latitude.FirewallUpdateAttributes{
				Name: "updated-example-firewall",
				Rules: []latitude.FirewallRule{
					{
						From:     "192.168.1.0/24",
						To:       "ANY",
						Port:     "80",
						Protocol: "TCP",
					},
					{
						From:     "10.0.0.0/8",
						To:       "ANY",
						Port:     "22",
						Protocol: "TCP",
					},
				},
			},
		},
	}

	firewall, _, err := client.Firewalls.Update(firewallID, updateRequest)
	return firewall, err
}

// createFirewallAssignment creates a new firewall assignment for a given server
func createFirewallAssignment(client *latitude.Client, firewallID, serverID string) (*latitude.FirewallAssignment, error) {
	createRequest := &latitude.FirewallAssignmentCreateRequest{
		Data: latitude.FirewallAssignmentCreateData{
			Type: "firewall_server",
			Attributes: latitude.FirewallAssignmentCreateAttributes{
				// This will be sent as "server_id" in the JSON due to the struct tag
				Server: serverID,
			},
		},
	}

	assignment, _, err := client.Firewalls.CreateAssignment(firewallID, createRequest)
	return assignment, err
}
