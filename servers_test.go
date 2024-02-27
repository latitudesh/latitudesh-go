package latitude

import (
	"testing"
)

const (
	testServerType = "servers"
)

func deleteServer(t *testing.T, c *Client, id string) {
	if _, err := c.Servers.Delete(id); err != nil {
		t.Fatal(err)
	}
}

func TestAccServerBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	// Create a new project
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
	defer deleteServer(t, c, s.ID)

	// Update newly created server
	rs := randString8()
	sur := ServerUpdateRequest{
		Data: ServerUpdateData{
			ID:   s.ID,
			Type: "servers",
			Attributes: ServerUpdateAttributes{
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
