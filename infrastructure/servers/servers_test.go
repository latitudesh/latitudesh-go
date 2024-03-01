package servers_test

import (
	"testing"

	latitude "github.com/latitudesh/latitudesh-go"
	servers "github.com/latitudesh/latitudesh-go/servers"
)

const (
	testServerType = "servers"
)

func deleteServer(t *testing.T, c *latitude.Client, id string) {
	if _, err := c.Servers.Delete(id); err != nil {
		t.Fatal(err)
	}
}

func TestAccServerBasic(t *testing.T) {
	latitude.SkipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := latitude.SetupWithProject(t)
	defer teardown()

	// Create a new project
	hn := latitude.RandString8()
	scr := servers.ServerCreateRequest{
		Data: servers.ServerCreateData{
			Type: testServerType,
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

	// Update newly created server
	rs := latitude.RandString8()
	sur := servers.ServerUpdateRequest{
		Data: servers.ServerUpdateData{
			ID:   s.ID,
			Type: "servers",
			Attributes: servers.ServerUpdateAttributes{
				Hostname: rs,
			},
		},
	}
	s, _, err = c.Servers.Update(s.ID, &sur)
	if err != nil {
		t.Fatal(err)
	}
	if s.Hostname != rs {
		t.Fatalf("Expected the hostname of the updated server to be %s, not %s", rs, s.Hostname)
	}

	dl, _, err := c.Servers.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(dl) != 1 {
		t.Fatalf("Server List should contain exactly one server, was: %v", dl)
	}
}
