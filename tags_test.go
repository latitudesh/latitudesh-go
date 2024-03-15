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

	t.Run("Create Tags", func(t *testing.T) {
		rs := randString8()
		tcr := TagCreateRequest{
			Data: TagCreateData{
				Type: testTagsType,
				Attributes: TagCreateAttributes{
					Name:        rs,
					Description: "Test Tag",
					Color:       "#ffffff",
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

	t.Run("Update Tags", func(t *testing.T) {
		rs := randString8()
		tur := TagUpdateRequest{
			Data: TagUpdateData{
				ID:   tagID,
				Type: testTagsType,
				Attributes: TagUpdateAttributes{
					Name:        rs,
					Description: "updated tag",
					Color:       "#fafadc",
				},
			},
		}

		tag, _, err := c.Tags.Update(tagID, &tur)
		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, tag.Name, rs, "Project Name")
	})

	t.Run("Get and List Tags", func(t *testing.T) {
		tagTest, _, err := c.Tags.Get(tagID)
		if err != nil {
			t.Fatal(err)
		}

		dl, _, err := c.Tags.List(nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(dl) < 1 {
			t.Fatal("There should be at least one tag created")
		}

		for _, tag := range dl {
			if tag.ID != tagTest.ID {
				continue
			}

			assertEqual(t, tag.Name, tagTest.Name, "Tag Name")
			assertEqual(t, tag.Slug, tagTest.Slug, "Tag Slug")
			assertEqual(t, tag.Description, tagTest.Description, "Tag Description")
			assertEqual(t, tag.Color, tagTest.Color, "Tag Color")
			assertEqual(t, tag.TeamID, tagTest.TeamID, "Tag TeamID")
			return
		}
		t.Fatalf("Tag with id %s not found", tagTest.ID)
	})
}
