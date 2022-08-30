package latitude

import (
	"path"
)

const serverBasePath = "/servers"

type ServerService interface {
	List(ProjectID string, opts *ListOptions) ([]ServerData, *Response, error)
	Get(ServerID string, opts *GetOptions) (*ServerGetResponse, *Response, error)
	Create(*ServerCreateRequest) (*Server, *Response, error)
	Update(string, *ServerUpdateRequest) (*Server, *Response, error)
	Delete(serverID string, force bool) (*Response, error)
}

type Server struct {
	Data ServerData `json:"data"`
	Meta meta       `json:"meta"`
}

type ServerData struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Attributes ServerAttributes `json:"attributes"`
}

type ServerAttributes struct {
	Hostname string      `json:"hostname"`
	Label    string      `json:"label"`
	Role     string      `json:"role"`
	Status   string      `json:"status"`
	Specs    ServerSpecs `json:"specs"`
}

type ServerSpecs struct {
	CPU  string `json:"cpu"`
	Disk string `json:"disk"`
	RAM  string `json:"ram"`
	NIC  string `json:"nic"`
}

type ServerListResponse struct {
	Servers []ServerData `json:"data"`
	Meta    meta         `json:"meta"`
}

type ServerGetResponse struct {
	Data ServerGetData `json:"data"`
	Meta meta          `json:"meta"`
}

type ServerGetData struct {
	ID         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes ServerGetAttributes `json:"attributes"`
}

type ServerGetAttributes struct {
	Hostname    string      `json:"hostname"`
	Label       string      `json:"label"`
	Role        string      `json:"role"`
	PrimaryIPv4 string      `json:"primary_ipv4"`
	Status      string      `json:"status"`
	IMPIStatus  string      `json:"impi_status"`
	CreatedAt   string      `json:"created_at"`
	Specs       ServerSpecs `json:"specs"`
}

// ServerCreateRequest type used to create a Latitude server
type ServerCreateRequest struct {
	Data ServerCreateData `json:"data"`
}

func (s ServerCreateRequest) String() string {
	return Stringify(s)
}

type ServerCreateData struct {
	Type       string                 `json:"type"`
	Attributes ServerCreateAttributes `json:"attributes"`
}

type ServerCreateAttributes struct {
	Project         string `json:"project,omitempty"`
	Plan            string `json:"plan,omitempty"`
	Site            string `json:"site,omitempty"`
	OperatingSystem string `json:"operating_system,omitempty"`
	Hostname        string `json:"hostname"`
}

// ServerUpdateRequest type used to update a Latitude server
type ServerUpdateRequest struct {
	Data ServerUpdateData `json:"data"`
}

type ServerUpdateData struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes ServerCreateAttributes `json:"attributes"`
}

func (p ServerUpdateRequest) String() string {
	return Stringify(p)
}

// ServerServiceOp implements ServerService
type ServerServiceOp struct {
	client requestDoer
}

// List returns servers on a project
func (s *ServerServiceOp) List(projectID string, opts *ListOptions) (servers []ServerData, resp *Response, err error) {
	opts = opts.Including("plan")
	endpointPath := path.Join(projectBasePath, projectID, serverBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		subset := new(ServerListResponse)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		servers = append(servers, subset.Servers...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}

// Get returns a server by id
func (s *ServerServiceOp) Get(serverID string, opts *GetOptions) (*ServerGetResponse, *Response, error) {
	endpointPath := path.Join(serverBasePath, serverID)
	apiPathQuery := opts.WithQuery(endpointPath)
	server := new(ServerGetResponse)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, server)
	if err != nil {
		return nil, resp, err
	}
	return server, resp, err
}

// Create creates a new server
func (s *ServerServiceOp) Create(createRequest *ServerCreateRequest) (*Server, *Response, error) {
	server := new(Server)

	resp, err := s.client.DoRequest("POST", serverBasePath, createRequest, server)
	if err != nil {
		return nil, resp, err
	}

	return server, resp, err
}

// Update updates a server
func (s *ServerServiceOp) Update(serverID string, updateRequest *ServerUpdateRequest) (*Server, *Response, error) {
	apiPath := path.Join(projectBasePath, serverID)
	server := new(Server)

	resp, err := s.client.DoRequest("PATCH", apiPath, updateRequest, server)
	if err != nil {
		return nil, resp, err
	}

	return server, resp, err
}

// Delete deletes a server
func (s *ServerServiceOp) Delete(serverID string, force bool) (*Response, error) {
	apiPath := path.Join(serverBasePath, serverID)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}
