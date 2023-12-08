package latitude

import (
	"fmt"
	"testing"
)

const (
	testProjectType        = "projects"
	testProjectEnvironment = "Development"
)

func TestAccProjectBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, stopRecord := setup(t)
	defer stopRecord()
	defer projectTeardown(c)

	// List Projects
	projs, _, err := c.Projects.List(nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, proj := range projs {
		fmt.Println(proj.ID, proj.Name)
	}

	// Create a new project
	rs := testProjectPrefix + randString8()
	pcr := ProjectCreateRequest{
		Data: ProjectCreateData{
			Type: testProjectType,
			Attributes: ProjectCreateAttributes{
				Name:        rs,
				Environment: testProjectEnvironment,
			},
		},
	}
	p, _, err := c.Projects.Create(&pcr)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != rs {
		t.Fatalf("Expected new project name to be %s, not %s", rs, p.Name)
	}

	// Update newly created project
	rs = testProjectPrefix + randString8()
	pur := ProjectUpdateRequest{
		Data: ProjectUpdateData{
			ID:   p.ID,
			Type: testProjectType,
			Attributes: ProjectCreateAttributes{
				Name:        rs,
				Environment: testProjectEnvironment,
			},
		},
	}
	p, _, err = c.Projects.Update(p.ID, &pur)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != rs {
		t.Fatalf("Expected the name of the updated project to be %s, not %s", rs, p.Name)
	}

	// Get newly updated project
	gotProject, _, err := c.Projects.Get(p.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if gotProject.Name != rs {
		t.Fatalf("Expected the name of the GOT project to be %s, not %s", rs, gotProject.Name)
	}

	// Delete newly created project
	_, err = c.Projects.Delete(p.ID)
	if err != nil {
		t.Fatal(err)
	}
}
