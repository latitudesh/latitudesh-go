package latitude

import (
	"testing"
)

const (
	testTagsType = "tags"
)

func deleteTag(t *testing.T, c *Client, id string) {
	if _, err := c.Tags.Delete(id); err != nil {
		t.Fatal(err)
	}
}

func TestAccTagBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, stopRecord := setup(t)
	defer stopRecord()

	var tagID string

	t.Run("Tags Create test", func(t *testing.T) {
		rs := randString8()
		tcr := TagCreateRequest{
			Data: TagCreateData{
				Type: testTagsType,
				Attributes: TagCreateAttributes{
					Name:        rs,
					Description: "Test Tag",
					Color:       "",
				},
			},
		}
		tag, _, err := c.Tags.Create(&tcr)
		if err != nil {
			t.Fatal(err)
		}

		tagID = tag.ID

		assertEqual(t, tag.Name, rs, "Tag Name")
	})

	// delete the tag at the end of the tests
	defer deleteTag(t, c, tagID)

	t.Run("Tags Update test", func(t *testing.T) {
	})

	t.Run("Tags List test", func(t *testing.T) {
		dl, _, err := c.Tags.List(nil)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, len(dl), 1, "Tag List length")
	})
}
