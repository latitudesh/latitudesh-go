package latitude

import (
	"errors"
	"net/http"
	"path"
)

const vlanAssignmentBasePath = "/virtual_networks/assignments"

type VlanAssignmentService interface {
	List(listOpt *ListOptions) ([]VlanAssignment, *Response, error)
	Get(VlanAssignmentID string) (*VlanAssignment, *Response, error)
	Assign(assignRequest *VlanAssignRequest) (*VlanAssignment, *Response, error)
	Delete(VlanAssignmentID string) (*Response, error)
}

type VlanAssignmentServiceOp struct {
	client requestDoer
}

type VlanAssignmentData struct {
	ID         string                   `json:"id"`
	Type       string                   `json:"type"`
	Attributes VlanAssignmentAttributes `json:"attributes"`
}

type VlanAssignmentAttributes struct {
	VlanAssignmentId string     `json:"virtual_network_id"`
	Vid              int        `json:"vid"`
	Description      string     `json:"description"`
	Status           string     `json:"status"`
	Server           VlanServer `json:"server"`
}

type VlanServer struct {
	Id       string `json:"id"`
	Hostname string `json:"hostname"`
	Label    string `json:"label"`
	Status   string `json:"status"`
}

type VlanAssignment struct {
	ID               string `json:"id"`
	Type             string `json:"type"`
	VirtualNetworkID string `json:"virtual_network_id"`
	Vid              int    `json:"vid"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	ServerID         string `json:"server_id"`
	ServerHostname   string `json:"server_hostname"`
	ServerStatus     string `json:"server_status"`
	ServerLabel      string `json:"server_label"`
}

type VlanAssignmentListResponse struct {
	Data []VlanAssignmentData `json:"data"`
	Meta meta                 `json:"meta"`
}

type VlanAssignmentGetResponse struct {
	Data VlanAssignmentData `json:"data"`
	Meta meta               `json:"meta"`
}

type VlanAssignmentCreateResponse struct {
	Data VlanAssignmentCreateData `json:"data"`
	Meta meta                     `json:"meta"`
}

type VlanAssignmentCreateData struct {
	ID         string                         `json:"id"`
	Type       string                         `json:"type"`
	Attributes VlanAssignmentCreateAttributes `json:"attributes"`
}

type VlanAssignmentCreateAttributes struct {
	VirtualNetworkID string `json:"virtual_network_id"`
	Vid              int    `json:"vid"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	ServerId         string `json:"server_id"`
}

type VlanAssignRequest struct {
	Data VlanAssignData `json:"data"`
}

type VlanAssignData struct {
	Type       string               `json:"type"`
	Attributes VlanAssignAttributes `json:"attributes"`
}

type VlanAssignAttributes struct {
	ServerID         string `json:"server_id"`
	VirtualNetworkID string `json:"virtual_network_id"`
}

func NewFlatVlanAssignment(vnd VlanAssignmentData) VlanAssignment {
	return VlanAssignment{
		vnd.ID,
		vnd.Type,
		vnd.Attributes.VlanAssignmentId,
		vnd.Attributes.Vid,
		vnd.Attributes.Description,
		vnd.Attributes.Status,
		vnd.Attributes.Server.Id,
		vnd.Attributes.Server.Hostname,
		vnd.Attributes.Server.Status,
		vnd.Attributes.Server.Label,
	}
}

func NewCreateFlatVlanAssignment(vnd VlanAssignmentCreateData) VlanAssignment {
	return VlanAssignment{
		ID:               vnd.ID,
		Type:             vnd.Type,
		VirtualNetworkID: vnd.Attributes.VirtualNetworkID,
		Vid:              vnd.Attributes.Vid,
		Description:      vnd.Attributes.Description,
		Status:           vnd.Attributes.Status,
		ServerID:         vnd.Attributes.ServerId,
	}
}

func NewFlatVlanAssignmentList(vnd []VlanAssignmentData) []VlanAssignment {
	var res []VlanAssignment
	for _, vn := range vnd {
		res = append(res, NewFlatVlanAssignment(vn))
	}
	return res
}

func (vn *VlanAssignmentServiceOp) List(opts *ListOptions) (vlanAssignments []VlanAssignment, resp *Response, err error) {
	apiPathQuery := opts.WithQuery(vlanAssignmentBasePath)

	for {
		res := new(VlanAssignmentListResponse)

		resp, err = vn.client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		vlanAssignments = append(vlanAssignments, NewFlatVlanAssignmentList(res.Data)...)

		if apiPathQuery = nextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}
}

func (s *VlanAssignmentServiceOp) Get(vlanAssignmentID string) (*VlanAssignment, *Response, error) {
	vlans, resp, err := s.List(nil)
	if err != nil {
		return nil, resp, err
	}

	for _, vlan := range vlans {
		if vlan.ID == vlanAssignmentID {
			return &vlan, resp, nil
		}
	}

	resp.Status = "404 Not Found"
	resp.StatusCode = http.StatusNotFound

	notFoundErr := errors.New("ERROR\nStatus: 404\nSpecified Record Not Found")

	return nil, resp, notFoundErr
}

func (s *VlanAssignmentServiceOp) Assign(assignRequest *VlanAssignRequest) (*VlanAssignment, *Response, error) {
	vLan := new(VlanAssignmentCreateResponse)

	resp, err := s.client.DoRequest("POST", vlanAssignmentBasePath, assignRequest, vLan)
	if err != nil {
		return nil, resp, err
	}

	flatVlanAssignment := NewCreateFlatVlanAssignment(vLan.Data)
	return &flatVlanAssignment, resp, err
}

func (s *VlanAssignmentServiceOp) Delete(vlanAssignmentID string) (*Response, error) {
	apiPath := path.Join(vlanAssignmentBasePath, vlanAssignmentID)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}
