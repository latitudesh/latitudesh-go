package latitude

import (
	"fmt"
	"path"
	"time"
)

const serverBasePath = "/servers"

type ServerService interface {
	List(ProjectID string, opts *ListOptions) ([]Server, *Response, error)
	Get(ServerID string, opts *GetOptions) (*Server, *Response, error)
	Create(*ServerCreateRequest) (*Server, *Response, error)
	Update(string, *ServerUpdateRequest) (*Server, *Response, error)
	Delete(serverID string) (*Response, error)
	Reinstall(serverID string, reinstallRequest *ServerReinstallRequest) (*Response, error)
	Lock(serverID string) (*Response, error)
	Unlock(serverID string) (*Response, error)
}

type ServerRoot struct {
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
	GPU  string `json:"gpu"`
}

type ServerTeam struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Slug        string       `json:"slug"`
	Description string       `json:"description"`
	Address     string       `json:"address"`
	Status      string       `json:"status"`
	Currency    TeamCurrency `json:"currency"`
}

type TeamCurrency struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type ServerListResponse struct {
	Data []ServerGetData `json:"data"`
	Meta meta            `json:"meta"`
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
	Hostname        string                `json:"hostname"`
	Label           string                `json:"label"`
	Price           float64               `json:"price"`
	Role            string                `json:"role"`
	PrimaryIPv4     string                `json:"primary_ipv4"`
	Status          string                `json:"status"`
	IMPIStatus      string                `json:"impi_status"`
	Site            string                `json:"site"`
	InstanceType    string                `json:"instance_type"`
	Locked          bool                  `json:"locked"`
	CreatedAt       string                `json:"created_at"`
	Specs           ServerSpecs           `json:"specs"`
	Project         ServerProject         `json:"project"`
	OperatingSystem ServerOperatingSystem `json:"operating_system"`
	Plan            ServerPlan            `json:"plan"`
	Region          ServerRegion          `json:"region"`
	Team            ServerTeam            `json:"team"`
	Tags            []EmbedTag            `json:"tags"`
}

// ServerCreateRequest type used to create a Latitude server
type ServerCreateRequest struct {
	Data ServerCreateData `json:"data"`
}

type ServerCreateData struct {
	Type       string                 `json:"type"`
	Attributes ServerCreateAttributes `json:"attributes"`
}

type ServerCreateAttributes struct {
	Project         string   `json:"project,omitempty"`
	Plan            string   `json:"plan,omitempty"`
	Site            string   `json:"site,omitempty"`
	OperatingSystem string   `json:"operating_system,omitempty"`
	Hostname        string   `json:"hostname"`
	SSHKeys         []string `json:"ssh_keys,omitempty"`
	UserData        string   `json:"user_data,omitempty"`
	Raid            string   `json:"raid,omitempty"`
	IpxeUrl         string   `json:"ipxe_url,omitempty"`
	Billing         string   `json:"billing,omitempty"`
}

// ServerUpdateRequest type used to update a Latitude server
type ServerUpdateRequest struct {
	Data ServerUpdateData `json:"data"`
}

type ServerUpdateData struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes ServerUpdateAttributes `json:"attributes"`
}

type ServerUpdateAttributes struct {
	Hostname string   `json:"hostname"`
	Billing  string   `json:"billing,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}
type ServerReinstallRequest struct {
	Data ServerReinstallData `json:"data"`
}

type ServerReinstallData struct {
	Type       string                    `json:"type"`
	Attributes ServerReinstallAttributes `json:"attributes"`
}

type ServerReinstallAttributes struct {
	OperatingSystem string   `json:"operating_system,omitempty"`
	Hostname        string   `json:"hostname"`
	SSHKeys         []string `json:"ssh_keys,omitempty"`
	UserData        string   `json:"user_data,omitempty"`
	Raid            string   `json:"raid,omitempty"`
	IpxeUrl         string   `json:"ipxe_url,omitempty"`
}

// ServerServiceOp implements ServerService
type ServerServiceOp struct {
	client requestDoer
}

type Server struct {
	ID              string                `json:"id"`
	Hostname        string                `json:"hostname"`
	Label           string                `json:"label"`
	Role            string                `json:"role"`
	Status          string                `json:"status"`
	PrimaryIPv4     string                `json:"primary_ipv4"`
	IMPIStatus      string                `json:"impi_status"`
	Locked          bool                  `json:"locked"`
	CreatedAt       string                `json:"created_at"`
	Specs           ServerSpecs           `json:"specs"`
	Project         ServerProject         `json:"project"`
	OperatingSystem ServerOperatingSystem `json:"operating_system"`
	Plan            ServerPlan            `json:"plan"`
	Region          ServerRegion          `json:"region"`
	Tags            []EmbedTag            `json:"tags"`
}

type ServerProject struct {
	ID   interface{} `json:"id"`
	Name string      `json:"name"`
}

type ServerRegion struct {
	City    string     `json:"city"`
	Country string     `json:"country"`
	Site    ServerSite `json:"site"`
}

type ServerSite struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Facility string `json:"facility"`
}

type ServerPlan struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ServerOperatingSystem struct {
	Name     string                  `json:"name"`
	Slug     string                  `json:"slug"`
	Version  string                  `json:"version"`
	Features OperatingSystemFeatures `json:"features"`
	Distro   OperatingSystemDistro   `json:"distro"`
}

// Flatten latitude API data structures
func NewFlatServer(sd ServerGetData) Server {
	return Server{
		sd.ID,
		sd.Attributes.Hostname,
		sd.Attributes.Label,
		sd.Attributes.Role,
		sd.Attributes.Status,
		sd.Attributes.PrimaryIPv4,
		sd.Attributes.IMPIStatus,
		sd.Attributes.Locked,
		sd.Attributes.CreatedAt,
		sd.Attributes.Specs,
		sd.Attributes.Project,
		sd.Attributes.OperatingSystem,
		sd.Attributes.Plan,
		sd.Attributes.Region,
		sd.Attributes.Tags,
	}
}

func NewFlatServerList(sd []ServerGetData) []Server {
	var res []Server
	for _, server := range sd {
		res = append(res, NewFlatServer(server))
	}
	return res
}

func waitServerActive(s *ServerServiceOp, id string) (*Server, error) {
	// 15 minutes = 180 * 15sec-retry
	for i := 0; i < 180; i++ {
		<-time.After(15 * time.Second)
		s, _, err := s.Get(id, nil)
		if err != nil {
			return nil, err
		}
		if s.Status == "on" {
			return s, nil
		}
		if s.Status == "failed" {
			return nil, fmt.Errorf("device %s provisioning failed", id)
		}
	}

	return nil, fmt.Errorf("device %s is still not active after timeout", id)
}

// List returns servers on a project
func (s *ServerServiceOp) List(projectID string, opts *ListOptions) ([]Server, *Response, error) {
	opts = opts.Filter("project", projectID)
	apiPathQuery := opts.WithQuery(serverBasePath)
	var servers []Server

	for {
		res := new(ServerListResponse)

		resp, err := s.client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		servers = append(servers, NewFlatServerList(res.Data)...)

		if apiPathQuery = nextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}

		return servers, resp, err
	}
}

// Get returns a server by id
func (s *ServerServiceOp) Get(serverID string, opts *GetOptions) (*Server, *Response, error) {
	endpointPath := path.Join(serverBasePath, serverID)
	apiPathQuery := opts.WithQuery(endpointPath)
	server := new(ServerGetResponse)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, server)
	if err != nil {
		return nil, resp, err
	}

	flatServer := NewFlatServer(server.Data)
	return &flatServer, resp, err
}

// Create creates a new server
func (s *ServerServiceOp) Create(createRequest *ServerCreateRequest) (*Server, *Response, error) {
	server := new(ServerGetResponse)

	resp, err := s.client.DoRequest("POST", serverBasePath, createRequest, server)
	if err != nil {
		return nil, resp, err
	}

	flatServer := NewFlatServer(server.Data)
	_, err = waitServerActive(s, flatServer.ID)
	return &flatServer, resp, err
}

// Update updates a server
func (s *ServerServiceOp) Update(serverID string, updateRequest *ServerUpdateRequest) (*Server, *Response, error) {
	apiPath := path.Join(serverBasePath, serverID)
	server := new(ServerGetResponse)

	resp, err := s.client.DoRequest("PATCH", apiPath, updateRequest, server)
	if err != nil {
		return nil, resp, err
	}

	flatServer := NewFlatServer(server.Data)
	return &flatServer, resp, err
}

// Delete deletes a server
func (s *ServerServiceOp) Delete(serverID string) (*Response, error) {
	apiPath := path.Join(serverBasePath, serverID)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}

// Reinstall reinstalls an existing server
func (s *ServerServiceOp) Reinstall(serverID string, reinstallRequest *ServerReinstallRequest) (*Response, error) {
	apiPath := path.Join(serverBasePath, serverID, "reinstall")

	return s.client.DoRequest("POST", apiPath, reinstallRequest, nil)
}

// Lock locks the server. A locked server cannot be deleted or modified and no actions can be performed on it.
func (s *ServerServiceOp) Lock(serverID string) (*Response, error) {
	apiPath := path.Join(serverBasePath, serverID, "lock")

	return s.client.DoRequest("POST", apiPath, nil, nil)
}

func (s *ServerServiceOp) Unlock(serverID string) (*Response, error) {
	apiPath := path.Join(serverBasePath, serverID, "unlock")

	return s.client.DoRequest("POST", apiPath, nil, nil)
}
