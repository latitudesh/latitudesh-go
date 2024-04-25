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

	tagIDs, deleteTags := setupTestTags(t, c)
	defer deleteTags()

	var serverId string

	t.Run("Servers Create test", func(t *testing.T) {
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
		serverId = s.ID
	})

	// delete the server at the end of the tests
	defer deleteServer(t, c, serverId)

	t.Run("Servers Update test", func(t *testing.T) {
		rs := randString8()
		sur := ServerUpdateRequest{
			Data: ServerUpdateData{
				ID:   serverId,
				Type: "servers",
				Attributes: ServerUpdateAttributes{
					Hostname: rs,
					Tags:     tagIDs,
				},
			},
		}
		s, _, err := c.Servers.Update(serverId, &sur)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, s.Hostname, rs, "Server hostname")
		assertEqual(t, len(s.Tags), 2, "Virtual Network Tags")
	})

	t.Run("Servers List test", func(t *testing.T) {
		dl, _, err := c.Servers.List(projectID, nil)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, len(dl), 1, "Server List length")
	})
}
