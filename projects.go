package latitude

import (
	"path"
)

const projectBasePath = "/projects"

// ProjectService interface defines available project methods
type ProjectService interface {
	List(listOpt *ListOptions) ([]Project, *Response, error)
	Get(string, *GetOptions) (*Project, *Response, error)
	Create(*ProjectCreateRequest) (*Project, *Response, error)
	Update(string, *ProjectUpdateRequest) (*Project, *Response, error)
	Delete(string) (*Response, error)
}

type ProjectRoot struct {
	Data ProjectData `json:"data"`
	Meta meta        `json:"meta"`
}

type ProjectData struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Attributes ProjectGetAttributes `json:"attributes"`
}

type ProjectListResponse struct {
	Data []ProjectData `json:"data"`
	Meta meta          `json:"meta"`
}

type ProjectGetResponse struct {
	Data ProjectData `json:"data"`
	Meta meta        `json:"meta"`
}

type ProjectGetAttributes struct {
	Name             string     `json:"name"`
	Slug             string     `json:"slug"`
	Description      string     `json:"description"`
	BillingType      string     `json:"billing_type"`
	BillingMethod    string     `json:"billing_method"`
	ProvisioningType string     `json:"provisioning_type"`
	Environment      string     `json:"environment"`
	CreatedAt        string     `json:"created_at"`
	UpdatedAt        string     `json:"updated_at"`
	Tags             []EmbedTag `json:"tags"`
}

// ProjectCreateRequest type used to create a Latitude project
type ProjectCreateRequest struct {
	Data ProjectCreateData `json:"data"`
}

type ProjectCreateData struct {
	Type       string                  `json:"type"`
	Attributes ProjectCreateAttributes `json:"attributes"`
}

type ProjectCreateAttributes struct {
	Name             string `json:"name"`
	ProvisioningType string `json:"provisioning_type"`
	Description      string `json:"description,omitempty"`
	Environment      string `json:"environment"`
}

// ProjectUpdateRequest type used to update a Latitude project
type ProjectUpdateRequest struct {
	Data ProjectUpdateData `json:"data"`
}

type ProjectUpdateData struct {
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Attributes ProjectUpdateAttributes `json:"attributes"`
}

type ProjectUpdateAttributes struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Environment string   `json:"environment"`
	Tags        []string `json:"tags,omitempty"`
}

// ProjectServiceOp implements ProjectService
type ProjectServiceOp struct {
	client requestDoer
}

type Project struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	Slug             string     `json:"slug"`
	Description      string     `json:"description"`
	BillingType      string     `json:"billing_type"`
	BillingMethod    string     `json:"billing_method"`
	ProvisioningType string     `json:"provisioning_type"`
	Environment      string     `json:"environment"`
	CreatedAt        string     `json:"created_at"`
	UpdatedAt        string     `json:"updated_at"`
	Tags             []EmbedTag `json:"tags"`
}

// Flatten latitude API data structures
func NewFlatProject(pd ProjectData) Project {
	return Project{
		pd.ID,
		pd.Attributes.Name,
		pd.Attributes.Slug,
		pd.Attributes.Description,
		pd.Attributes.BillingType,
		pd.Attributes.BillingMethod,
		pd.Attributes.ProvisioningType,
		pd.Attributes.Environment,
		pd.Attributes.CreatedAt,
		pd.Attributes.UpdatedAt,
		pd.Attributes.Tags,
	}
}

func NewFlatProjectList(pd []ProjectData) []Project {
	var res []Project
	for _, project := range pd {
		res = append(res, NewFlatProject(project))
	}
	return res
}

// List returns a list of projects
func (s *ProjectServiceOp) List(opts *ListOptions) (projects []Project, resp *Response, err error) {
	apiPathQuery := opts.WithQuery(projectBasePath)

	for {
		res := new(ProjectListResponse)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		projects = append(projects, NewFlatProjectList(res.Data)...)

		if apiPathQuery = nextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}

// Get returns a project by id
func (s *ProjectServiceOp) Get(projectID string, opts *GetOptions) (*Project, *Response, error) {
	endpointPath := path.Join(projectBasePath, projectID)
	apiPathQuery := opts.WithQuery(endpointPath)
	project := new(ProjectGetResponse)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, project)
	if err != nil {
		return nil, resp, err
	}

	flatProject := NewFlatProject(project.Data)
	return &flatProject, resp, err
}

// Create creates a new project
func (s *ProjectServiceOp) Create(createRequest *ProjectCreateRequest) (*Project, *Response, error) {
	project := new(ProjectGetResponse)

	if createRequest.Data.Attributes.ProvisioningType == "" {
		createRequest.Data.Attributes.ProvisioningType = "reserved"
	}

	resp, err := s.client.DoRequest("POST", projectBasePath, createRequest, project)
	if err != nil {
		return nil, resp, err
	}

	flatProject := NewFlatProject(project.Data)
	return &flatProject, resp, err
}

// Update updates a project
func (s *ProjectServiceOp) Update(projectID string, updateRequest *ProjectUpdateRequest) (*Project, *Response, error) {
	apiPath := path.Join(projectBasePath, projectID)
	project := new(ProjectGetResponse)

	resp, err := s.client.DoRequest("PATCH", apiPath, updateRequest, project)
	if err != nil {
		return nil, resp, err
	}

	flatProject := NewFlatProject(project.Data)
	return &flatProject, resp, err
}

// Delete deletes a project
func (s *ProjectServiceOp) Delete(projectID string) (*Response, error) {
	apiPath := path.Join(projectBasePath, projectID)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}
