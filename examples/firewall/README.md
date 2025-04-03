# Latitude.sh Firewall SDK Example

This example demonstrates how to use the Latitude.sh Go SDK to work with Firewalls and Firewall Assignments.

## Prerequisites

- Go 1.16 or later
- A Latitude.sh account with API access
- An existing project in your Latitude.sh account

## Setup

1. Set the required environment variables:

```sh
export LATITUDE_API_TOKEN="your_api_token_here"
export LATITUDE_PROJECT="your_project_slug_here"
```

2. Run the example:

```sh
go run main.go
```

## What This Example Does

The example demonstrates the following operations:

1. Lists existing firewalls in your Latitude.sh account
2. Creates a new firewall with example rules
3. Updates the firewall with modified rules
4. Lists servers available for assignment (if any)
5. Creates a firewall assignment for a server (if available)
6. Lists firewall assignments
7. Deletes the firewall assignment
8. Deletes the firewall

## Example Output

```
Listing existing firewalls...
Found 1 firewalls
1. test-firewall (ID: fw_xkjQwdENqYNVP) - Project: My Project

Creating a new firewall...
Created firewall: example-firewall (ID: fw_VLMmAD8EOwop2)
Rules:
  1. From: 192.168.1.0/24, To: ANY, Port: 80, Protocol: TCP
  2. From: 10.0.0.0/8, To: 192.168.1.100, Port: 443, Protocol: TCP
  3. From: 172.16.0.0/12, To: ANY, Port: 3000-4000, Protocol: UDP

Updating the firewall...
Updated firewall: updated-example-firewall (ID: fw_VLMmAD8EOwop2)
Updated rules:
  1. From: 192.168.1.0/24, To: ANY, Port: 80, Protocol: TCP
  2. From: 10.0.0.0/8, To: ANY, Port: 22, Protocol: TCP

Listing servers for potential assignment...
Found 1 servers. First server: web-server (ID: srv_vG0EmKp7Ak3qX)

Creating firewall assignment...
Created firewall assignment (ID: fwsrv_6nkQ1EARx2jlz) for server: web-server

Listing firewall assignments...
Found 1 assignments

Deleting firewall assignment...
Successfully deleted firewall assignment

Deleting firewall...
Successfully deleted firewall

Firewall test completed successfully!
```

## Next Steps

After running this example, you can:

1. Modify the example code to create different firewall rules
2. Integrate firewall management into your own Go applications
3. Explore other features of the Latitude.sh SDK 