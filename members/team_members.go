package members

import (
	"path"

	api "github.com/latitudesh/latitudesh-go/api_utils"
	internal "github.com/latitudesh/latitudesh-go/internal"
	types "github.com/latitudesh/latitudesh-go/types"
)

const memberBasePath = "/team/members"

// MemberService interface defines available member methods
type MemberService interface {
	List(listOpt *api.ListOptions) ([]Member, *types.Response, error)
	Create(request *MemberCreateRequest) (*Member, *types.Response, error)
	Delete(UserID string) (*types.Response, error)
}

// MemberServiceOp implements MemberService
type MemberServiceOp struct {
	Client internal.RequestDoer
}

type MemberListResponse struct {
	Data []MemberListData `json:"data"`
	Meta internal.Meta    `json:"meta"`
}

type MemberListData struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Attributes MemberListAttributes `json:"attributes"`
}

type MemberListAttributes struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	MfaEnabled bool   `json:"mfa_enabled"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Role       Role   `json:"role"`
}

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type MemberResponse struct {
	Data MemberData    `json:"data"`
	Meta internal.Meta `json:"meta"`
}

type MemberData struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Attributes MemberAttributes `json:"attributes"`
}

type MemberAttributes struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	MfaEnabled bool   `json:"mfa_enabled"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	RoleName   string `json:"role"`
}

type MemberCreateRequest struct {
	Data MemberCreateData `json:"data"`
}

type MemberCreateData struct {
	Type       string                 `json:"type"`
	Attributes MemberCreateAttributes `json:"attributes"`
}

type MemberCreateAttributes struct {
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Role      MemberRole `json:"role"`
}

type MemberRole string

const (
	Owner         MemberRole = "owner"
	Administrator MemberRole = "administrator"
	Collaborator  MemberRole = "collaborator"
	Billing       MemberRole = "billing"
)

type Member struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	MfaEnabled bool   `json:"mfa_enabled"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	RoleName   string `json:"role"`
}

// Flatten latitude API data structure
func NewFlatMember(md MemberData) Member {
	return Member{
		md.ID,
		md.Attributes.FirstName,
		md.Attributes.LastName,
		md.Attributes.Email,
		md.Attributes.MfaEnabled,
		md.Attributes.CreatedAt,
		md.Attributes.UpdatedAt,
		md.Attributes.RoleName,
	}
}

func NewFlatMemberList(md []MemberData) []Member {
	var members []Member
	for _, member := range md {
		members = append(members, NewFlatMember(member))
	}
	return members
}

// List returns a list of team members
func (s *MemberServiceOp) List(listOpts *api.ListOptions) (members []Member, resp *types.Response, err error) {
	apiPathQuery := listOpts.WithQuery(memberBasePath)

	for {
		res := new(MemberListResponse)
		membersData := []MemberData{}

		resp, err = s.Client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		for _, data := range res.Data {
			mD := MemberData{
				ID:   data.ID,
				Type: data.Type,
				Attributes: MemberAttributes{
					FirstName:  data.Attributes.FirstName,
					LastName:   data.Attributes.LastName,
					Email:      data.Attributes.Email,
					MfaEnabled: data.Attributes.MfaEnabled,
					CreatedAt:  data.Attributes.CreatedAt,
					UpdatedAt:  data.Attributes.UpdatedAt,
					RoleName:   data.Attributes.Role.Name,
				},
			}

			membersData = append(membersData, mD)
		}

		members = append(members, NewFlatMemberList(membersData)...)

		if apiPathQuery = api.NextPage(res.Meta, listOpts); apiPathQuery != "" {
			continue
		}

		return
	}
}

// Create creates a new team member
func (s *MemberServiceOp) Create(request *MemberCreateRequest) (*Member, *types.Response, error) {
	member := new(MemberResponse)

	resp, err := s.Client.DoRequest("POST", memberBasePath, request, member)
	if err != nil {
		return nil, resp, err
	}

	flatMember := NewFlatMember(member.Data)
	return &flatMember, resp, err
}

// Delete deletes a team member
func (s *MemberServiceOp) Delete(MemberID string) (*types.Response, error) {
	apiPath := path.Join(memberBasePath, MemberID)

	return s.Client.DoRequest("DELETE", apiPath, nil, nil)
}
