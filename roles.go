package latitude

import "path"

const roleBasePath = "/roles"

// RoleService interface defines available role methods
type RoleService interface {
	Get(string, *GetOptions) (*Role, *Response, error)
	List(*ListOptions) ([]Role, *Response, error)
}

// RoleServiceOp implements RoleService
type RoleServiceOp struct {
	client requestDoer
}

type AvailableRole string

const (
	Owner         AvailableRole = "owner"
	Administrator AvailableRole = "administrator"
	Collaborator  AvailableRole = "collaborator"
	Billing       AvailableRole = "billing"
)

type RoleGetResponse struct {
	Data RoleData `json:"data"`
	Meta meta     `json:"meta"`
}

type RoleData struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes RoleAttributes `json:"attributes"`
}

type RoleAttributes struct {
	Name string `json:"name"`
}

type RoleListResponse struct {
	Data []RoleData `json:"data"`
	Meta meta       `json:"meta"`
}

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewFlatRole(rd RoleData) Role {
	return Role{
		ID:   rd.ID,
		Name: rd.Attributes.Name,
	}
}

func NewFlatRoleList(rd []RoleData) []Role {
	var res []Role
	for _, role := range rd {
		res = append(res, NewFlatRole(role))
	}
	return res
}

func (s *RoleServiceOp) Get(RoleID string, opts *GetOptions) (*Role, *Response, error) {
	endpointPath := path.Join(roleBasePath, RoleID)
	apiPathQuery := opts.WithQuery(endpointPath)
	role := new(RoleGetResponse)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, role)
	if err != nil {
		return nil, resp, err
	}

	flatRole := NewFlatRole(role.Data)
	return &flatRole, resp, err
}

func (s *RoleServiceOp) List(opts *ListOptions) ([]Role, *Response, error) {
	apiPathQuery := opts.WithQuery(roleBasePath)
	roles := []Role{}

	for {
		res := new(RoleListResponse)

		resp, err := s.client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		roles = append(roles, NewFlatRoleList(res.Data)...)

		if apiPathQuery = nextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}

		return roles, resp, nil
	}
}
