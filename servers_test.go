package latitude

import (
	"fmt"
	"testing"
	"time"
)

const (
	testServerType = "servers"
)

func waitServerActive(t *testing.T, c *Client, id string) *Server {
	// 15 minutes = 180 * 15sec-retry
	for i := 0; i < 180; i++ {
		<-time.After(15 * time.Second)
		d, _, err := c.Servers.Get(id, nil)
		if err != nil {
			t.Fatal(err)
			return nil
		}
		if d.Status == "on" {
			return d
		}
		if d.Status == "failed" {
			t.Fatal(fmt.Errorf("device %s provisioning failed", id))
			return nil
		}
	}

	t.Fatal(fmt.Errorf("device %s is still not active after timeout", id))
	return nil
}

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

	sID := s.ID

	sgr := waitServerActive(t, c, sID)

	if len(sgr.ID) == 0 {
		t.Fatal("Server should have an ID")
	}

	// TODO: API endpoint for server update currently not working
	// Update newly created server
	/*rs := randString8()
	sur := ServerUpdateRequest{
		Data: ServerUpdateData{
			ID:   sID,
			Type: testProjectType,
			Attributes: ServerCreateAttributes{
				Hostname: rs,
			},
		},
	}
	s, _, err = c.Servers.Update(sID, &sur)
	if err != nil {
		t.Fatal(err)
	}
	if s.Data.Attributes.Hostname != rs {
		t.Fatalf("Expected the hostname of the updated server to be %s, not %s", rs, s.Data.Attributes.Hostname)
	}*/

	dl, _, err := c.Servers.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(dl) != 1 {
		t.Fatalf("Server List should contain exactly one server, was: %v", dl)
	}
}
