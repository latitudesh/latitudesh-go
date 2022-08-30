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
	if p.Data.Attributes.Name != rs {
		t.Fatalf("Expected new project name to be %s, not %s", rs, p.Data.Attributes.Name)
	}

	// Update newly created project
	rs = testProjectPrefix + randString8()
	pur := ProjectUpdateRequest{
		Data: ProjectUpdateData{
			ID:   p.Data.ID,
			Type: testProjectType,
			Attributes: ProjectCreateAttributes{
				Name:        rs,
				Environment: testProjectEnvironment,
			},
		},
	}
	p, _, err = c.Projects.Update(p.Data.ID, &pur)
	if err != nil {
		t.Fatal(err)
	}
	if p.Data.Attributes.Name != rs {
		t.Fatalf("Expected the name of the updated project to be %s, not %s", rs, p.Data.Attributes.Name)
	}

	// Get newly updated project
	gotProject, _, err := c.Projects.Get(p.Data.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if gotProject.Data.Attributes.Name != rs {
		t.Fatalf("Expected the name of the GOT project to be %s, not %s", rs, gotProject.Data.Attributes.Name)
	}

	// Delete newly created project
	_, err = c.Projects.Delete(p.Data.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccListProjects(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, stopRecord := setup(t)
	defer stopRecord()

	listOpt := &ListOptions{
		Includes: []string{"team"},
	}
	projs, _, err := c.Projects.List(listOpt)
	if err != nil {
		t.Fatal(err)
	}

	for _, proj := range projs {
		fmt.Println(proj.ID, proj.Attributes.Name)
	}
}
