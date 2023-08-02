package latitude

import "path"

const vlanAssignmentBasePath = "/virtual_networks/assignments"

type VlanAssignmentService interface {
	List(listOpt *ListOptions) ([]VlanAssignment, *Response, error)
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
	VlanAssignmentId int        `json:"virtual_network_id"`
	Vid              int        `json:"vid"`
	Description      string     `json:"description"`
	Status           string     `json:"status"`
	Server           VlanServer `json:"server"`
}

type VlanServer struct {
	Id       int    `json:"id"`
	Hostname string `json:"hostname"`
	Label    string `json:"label"`
	Status   string `json:"status"`
}

type VlanAssignment struct {
	ID               string `json:"id"`
	Type             string `json:"type"`
	VlanAssignmentID int    `json:"virtual_network_id"`
	Vid              int    `json:"vid"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	ServerID         int    `json:"server_id"`
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

type VlanAssignmentCreateRequest struct {
	ServerID         int `json:"server_id"`
	VirtualNetworkID int `json:"virtual_network_id"`
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

func (s *VlanAssignmentServiceOp) Assign(createRequest *VlanAssignmentCreateRequest) (*VlanAssignment, *Response, error) {
	vLan := new(VlanAssignmentGetResponse)

	resp, err := s.client.DoRequest("POST", vlanAssignmentBasePath, createRequest, vLan)
	if err != nil {
		return nil, resp, err
	}

	flatVlanAssignment := NewFlatVlanAssignment(vLan.Data)
	return &flatVlanAssignment, resp, err
}

func (s *VlanAssignmentServiceOp) Delete(vlanAssignmentID string) (*Response, error) {
	apiPath := path.Join(vlanAssignmentBasePath, vlanAssignmentID)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}
