package virtual_network_assignment

import (
	"errors"
	"net/http"
	"path"

	api "github.com/latitudesh/latitudesh-go/api_utils"
	internal "github.com/latitudesh/latitudesh-go/internal"
	types "github.com/latitudesh/latitudesh-go/types"
)

const vlanAssignmentBasePath = "/virtual_networks/assignments"

type VlanAssignmentService interface {
	List(listOpt *api.ListOptions) ([]VlanAssignment, *types.Response, error)
	Get(VlanAssignmentID string) (*VlanAssignment, *types.Response, error)
	Assign(assignRequest *VlanAssignRequest) (*VlanAssignment, *types.Response, error)
	Delete(VlanAssignmentID string) (*types.Response, error)
}

type VlanAssignmentServiceOp struct {
	Client internal.RequestDoer
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
	Meta internal.Meta        `json:"meta"`
}

type VlanAssignmentGetResponse struct {
	Data VlanAssignmentData `json:"data"`
	Meta internal.Meta      `json:"meta"`
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

func NewFlatVlanAssignmentList(vnd []VlanAssignmentData) []VlanAssignment {
	var res []VlanAssignment
	for _, vn := range vnd {
		res = append(res, NewFlatVlanAssignment(vn))
	}
	return res
}

func (vn *VlanAssignmentServiceOp) List(opts *api.ListOptions) (vlanAssignments []VlanAssignment, resp *types.Response, err error) {
	apiPathQuery := opts.WithQuery(vlanAssignmentBasePath)

	for {
		res := new(VlanAssignmentListResponse)

		resp, err = vn.Client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		vlanAssignments = append(vlanAssignments, NewFlatVlanAssignmentList(res.Data)...)

		if apiPathQuery = api.NextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}
}

func (s *VlanAssignmentServiceOp) Get(vlanAssignmentID string) (*VlanAssignment, *types.Response, error) {
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

func (s *VlanAssignmentServiceOp) Assign(assignRequest *VlanAssignRequest) (*VlanAssignment, *types.Response, error) {
	vLan := new(VlanAssignmentGetResponse)

	resp, err := s.Client.DoRequest("POST", vlanAssignmentBasePath, assignRequest, vLan)
	if err != nil {
		return nil, resp, err
	}

	flatVlanAssignment := NewFlatVlanAssignment(vLan.Data)
	return &flatVlanAssignment, resp, err
}

func (s *VlanAssignmentServiceOp) Delete(vlanAssignmentID string) (*types.Response, error) {
	apiPath := path.Join(vlanAssignmentBasePath, vlanAssignmentID)

	return s.Client.DoRequest("DELETE", apiPath, nil, nil)
}
