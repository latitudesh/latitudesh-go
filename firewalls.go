package latitude

import (
	"path"
)

const firewallBasePath = "/firewalls"

// FirewallService interface defines available firewall methods
type FirewallService interface {
	List(listOpt *ListOptions) ([]Firewall, *Response, error)
	Get(string, *GetOptions) (*Firewall, *Response, error)
	Create(*FirewallCreateRequest) (*Firewall, *Response, error)
	Update(string, *FirewallUpdateRequest) (*Firewall, *Response, error)
	Delete(string) (*Response, error)
	ListAssignments(firewallID string, listOpt *ListOptions) ([]FirewallAssignment, *Response, error)
	CreateAssignment(firewallID string, request *FirewallAssignmentCreateRequest) (*FirewallAssignment, *Response, error)
	DeleteAssignment(firewallID string, assignmentID string) (*Response, error)
}

// FirewallRule represents a rule in a firewall
type FirewallRule struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	Default  bool   `json:"default"`
}

// FirewallProject represents embedded project information in firewall
type FirewallProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Firewall represents a Latitude Firewall
type Firewall struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Project FirewallProject `json:"project"`
	Rules   []FirewallRule  `json:"rules"`
}

// FirewallServiceOp implements FirewallService
type FirewallServiceOp struct {
	client requestDoer
}

// FirewallData represents the data structure returned by the API
type FirewallData struct {
	ID         string                `json:"id"`
	Type       string                `json:"type"`
	Attributes FirewallGetAttributes `json:"attributes"`
}

// FirewallGetAttributes represents the attributes in the API response
type FirewallGetAttributes struct {
	Name    string          `json:"name"`
	Project FirewallProject `json:"project"`
	Rules   []FirewallRule  `json:"rules"`
}

// FirewallListResponse represents the list response from the API
type FirewallListResponse struct {
	Data []FirewallData `json:"data"`
	Meta meta           `json:"meta"`
}

// FirewallGetResponse represents the get response from the API
type FirewallGetResponse struct {
	Data FirewallData `json:"data"`
	Meta meta         `json:"meta"`
}

// FirewallCreateRequest type used to create a Latitude Firewall
type FirewallCreateRequest struct {
	Data FirewallCreateData `json:"data"`
}

// FirewallCreateData represents the data structure for creating a firewall
type FirewallCreateData struct {
	Type       string                   `json:"type"`
	Attributes FirewallCreateAttributes `json:"attributes"`
}

// FirewallCreateAttributes represents the attributes for creating a firewall
type FirewallCreateAttributes struct {
	Name    string         `json:"name"`
	Project string         `json:"project"`
	Rules   []FirewallRule `json:"rules,omitempty"`
}

// FirewallUpdateRequest type used to update a Latitude Firewall
type FirewallUpdateRequest struct {
	Data FirewallUpdateData `json:"data"`
}

// FirewallUpdateData represents the data structure for updating a firewall
type FirewallUpdateData struct {
	ID         string                   `json:"id"`
	Type       string                   `json:"type"`
	Attributes FirewallUpdateAttributes `json:"attributes"`
}

// FirewallUpdateAttributes represents the attributes for updating a firewall
type FirewallUpdateAttributes struct {
	Name  string         `json:"name,omitempty"`
	Rules []FirewallRule `json:"rules,omitempty"`
}

// FirewallAssignment represents a firewall assignment to a server
type FirewallAssignment struct {
	ID     string `json:"id"`
	Server Server `json:"server"`
}

// FirewallAssignmentData represents the data structure returned by the API
type FirewallAssignmentData struct {
	ID         string                       `json:"id"`
	Type       string                       `json:"type"`
	Attributes FirewallAssignmentAttributes `json:"attributes"`
}

// FirewallAssignmentAttributes represents the attributes in the API response
type FirewallAssignmentAttributes struct {
	Server Server `json:"server"`
}

// FirewallAssignmentListResponse represents the list response from the API
type FirewallAssignmentListResponse struct {
	Data []FirewallAssignmentData `json:"data"`
	Meta meta                     `json:"meta"`
}

// FirewallAssignmentGetResponse represents the get response from the API
type FirewallAssignmentGetResponse struct {
	Data FirewallAssignmentData `json:"data"`
	Meta meta                   `json:"meta"`
}

// FirewallAssignmentCreateRequest type used to create a Latitude Firewall Assignment
type FirewallAssignmentCreateRequest struct {
	Data FirewallAssignmentCreateData `json:"data"`
}

// FirewallAssignmentCreateData represents the data structure for creating a firewall assignment
type FirewallAssignmentCreateData struct {
	Type       string                             `json:"type"`
	Attributes FirewallAssignmentCreateAttributes `json:"attributes"`
}

// FirewallAssignmentCreateAttributes represents the attributes for creating a firewall assignment
type FirewallAssignmentCreateAttributes struct {
	Server string `json:"server_id"`
}

// NewFlatFirewall flattens the API response to a Firewall struct
func NewFlatFirewall(fd FirewallData) Firewall {
	return Firewall{
		ID:      fd.ID,
		Name:    fd.Attributes.Name,
		Project: fd.Attributes.Project,
		Rules:   fd.Attributes.Rules,
	}
}

// NewFlatFirewallList flattens a list of API responses to a list of Firewall structs
func NewFlatFirewallList(fd []FirewallData) []Firewall {
	var firewalls []Firewall
	for _, firewall := range fd {
		firewalls = append(firewalls, NewFlatFirewall(firewall))
	}
	return firewalls
}

// NewFlatFirewallAssignment flattens the API response to a FirewallAssignment struct
func NewFlatFirewallAssignment(fd FirewallAssignmentData) FirewallAssignment {
	return FirewallAssignment{
		ID:     fd.ID,
		Server: fd.Attributes.Server,
	}
}

// NewFlatFirewallAssignmentList flattens a list of API responses to a list of FirewallAssignment structs
func NewFlatFirewallAssignmentList(fd []FirewallAssignmentData) []FirewallAssignment {
	var assignments []FirewallAssignment
	for _, assignment := range fd {
		assignments = append(assignments, NewFlatFirewallAssignment(assignment))
	}
	return assignments
}

// List returns a list of firewalls
func (s *FirewallServiceOp) List(opts *ListOptions) (firewalls []Firewall, resp *Response, err error) {
	apiPathQuery := opts.WithQuery(firewallBasePath)

	for {
		res := new(FirewallListResponse)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		firewalls = append(firewalls, NewFlatFirewallList(res.Data)...)

		if apiPathQuery = nextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}

// Get returns a firewall by id
func (s *FirewallServiceOp) Get(firewallID string, opts *GetOptions) (*Firewall, *Response, error) {
	endpointPath := path.Join(firewallBasePath, firewallID)
	apiPathQuery := opts.WithQuery(endpointPath)
	firewall := new(FirewallGetResponse)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, firewall)
	if err != nil {
		return nil, resp, err
	}

	flatFirewall := NewFlatFirewall(firewall.Data)
	return &flatFirewall, resp, err
}

// Create creates a new firewall
func (s *FirewallServiceOp) Create(createRequest *FirewallCreateRequest) (*Firewall, *Response, error) {
	firewall := new(FirewallGetResponse)

	// Set type if not specified
	if createRequest.Data.Type == "" {
		createRequest.Data.Type = "firewalls"
	}

	resp, err := s.client.DoRequest("POST", firewallBasePath, createRequest, firewall)
	if err != nil {
		return nil, resp, err
	}

	flatFirewall := NewFlatFirewall(firewall.Data)
	return &flatFirewall, resp, err
}

// Update updates a firewall
func (s *FirewallServiceOp) Update(firewallID string, updateRequest *FirewallUpdateRequest) (*Firewall, *Response, error) {
	apiPath := path.Join(firewallBasePath, firewallID)
	firewall := new(FirewallGetResponse)

	// Set type if not specified
	if updateRequest.Data.Type == "" {
		updateRequest.Data.Type = "firewalls"
	}

	resp, err := s.client.DoRequest("PATCH", apiPath, updateRequest, firewall)
	if err != nil {
		return nil, resp, err
	}

	flatFirewall := NewFlatFirewall(firewall.Data)
	return &flatFirewall, resp, err
}

// Delete deletes a firewall
func (s *FirewallServiceOp) Delete(firewallID string) (*Response, error) {
	apiPath := path.Join(firewallBasePath, firewallID)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}

// ListAssignments returns a list of firewall assignments
func (s *FirewallServiceOp) ListAssignments(firewallID string, opts *ListOptions) (assignments []FirewallAssignment, resp *Response, err error) {
	apiPath := path.Join(firewallBasePath, firewallID, "assignments")
	apiPathQuery := opts.WithQuery(apiPath)

	for {
		res := new(FirewallAssignmentListResponse)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		assignments = append(assignments, NewFlatFirewallAssignmentList(res.Data)...)

		if apiPathQuery = nextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}

// CreateAssignment creates a new firewall assignment
func (s *FirewallServiceOp) CreateAssignment(firewallID string, createRequest *FirewallAssignmentCreateRequest) (*FirewallAssignment, *Response, error) {
	apiPath := path.Join(firewallBasePath, firewallID, "assignments")
	assignment := new(FirewallAssignmentGetResponse)

	// Set type if not specified
	if createRequest.Data.Type == "" {
		createRequest.Data.Type = "firewall_server"
	}

	resp, err := s.client.DoRequest("POST", apiPath, createRequest, assignment)
	if err != nil {
		return nil, resp, err
	}

	flatAssignment := NewFlatFirewallAssignment(assignment.Data)
	return &flatAssignment, resp, err
}

// DeleteAssignment deletes a firewall assignment
func (s *FirewallServiceOp) DeleteAssignment(firewallID string, assignmentID string) (*Response, error) {
	apiPath := path.Join(firewallBasePath, firewallID, "assignments", assignmentID)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}
