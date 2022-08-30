package latitude

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

const (
	testServerType = "servers"
)

func waitServerActive(t *testing.T, c *Client, id string) *ServerGetResponse {
	// 15 minutes = 180 * 5sec-retry
	for i := 0; i < 180; i++ {
		<-time.After(5 * time.Second)
		d, _, err := c.Servers.Get(id, nil)
		if err != nil {
			t.Fatal(err)
			return nil
		}
		if d.Data.Attributes.Status == "on" {
			return d
		}
		if d.Data.Attributes.Status == "failed" {
			t.Fatal(fmt.Errorf("device %s provisioning failed", id))
			return nil
		}
	}

	t.Fatal(fmt.Errorf("device %s is still not active after timeout", id))
	return nil
}

func deleteServer(t *testing.T, c *Client, id string, force bool) {
	if _, err := c.Servers.Delete(id, force); err != nil {
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
				OperatingSystem: testOS,
				Hostname:        hn,
			},
		},
	}

	s, _, err := c.Servers.Create(&scr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteServer(t, c, s.Data.ID, false)

	sID := s.Data.ID

	sgr := waitServerActive(t, c, sID)

	if len(sgr.Data.ID) == 0 {
		t.Fatal("Server should have an ID")
	}

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
		t.Fatalf("Device List should contain exactly one device, was: %v", dl)
	}
}
