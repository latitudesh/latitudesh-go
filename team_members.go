package latitude

const memberBasePath = "/team/members"

// MemberService interface defines available member methods
type MemberService interface {
	List(listOpt *ListOptions) ([]User, *Response, error)
	Create(request *MemberCreateRequest) (*User, *Response, error)
	Delete(UserID string) (*Response, error)
}

// MemberServiceOp implements MemberService
type MemberServiceOp struct {
	client requestDoer
}

type MemberListResponse struct {
	Data []MemberData `json:"data"`
	Meta meta         `json:"meta"`
}

type MemberData struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Attributes MemberAttributes `json:"attributes"`
}

type MemberAttributes struct {
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Email      string     `json:"email"`
	MfaEnabled bool       `json:"mfa_enabled"`
	CreatedAt  string     `json:"created_at"`
	UpdatedAt  string     `json:"updated_at"`
	Role       RoleStruct `json:"role"`
}

type MemberCreateRequest struct {
	Data MemberCreateData `json:"data"`
	Meta meta             `json:"meta"`
}

type MemberCreateData struct {
	Type       string                 `json:"type"`
	Attributes MemberCreateAttributes `json:"attributes"`
}

type MemberCreateAttributes struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Role      UserRole `json:"role"`
}

type UserRole string

const (
	Owner         UserRole = "owner"
	Administrator UserRole = "administrator"
	Collaborator  UserRole = "collaborator"
	Billing       UserRole = "billing"
)

type User struct {
	ID         string     `json:"id"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Email      string     `json:"email"`
	MfaEnabled bool       `json:"mfa_enabled"`
	CreatedAt  string     `json:"created_at"`
	UpdatedAt  string     `json:"updated_at"`
	Role       RoleStruct `json:"role"`
}

type RoleStruct struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Flatten latitude API data structures
func NewFlatMember(md MemberData) User {
	return User{
		md.ID,
		md.Attributes.FirstName,
		md.Attributes.LastName,
		md.Attributes.Email,
		md.Attributes.MfaEnabled,
		md.Attributes.CreatedAt,
		md.Attributes.UpdatedAt,
		md.Attributes.Role,
	}
}

func NewFlatMemberList(md []MemberData) []User {
	var res []User
	for _, member := range md {
		res = append(res, NewFlatMember(member))
	}
	return res
}

func (s *MemberServiceOp) List(listOpts *ListOptions) (members []User, resp *Response, err error) {
	apiPathQuery := listOpts.WithQuery(memberBasePath)

	for {
		res := new(MemberListResponse)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		members = append(members, NewFlatMemberList(res.Data)...)

		if apiPathQuery = nextPage(res.Meta, listOpts); apiPathQuery != "" {
			continue
		}

		return
	}
}

func (s *MemberServiceOp) Create(request *MemberCreateRequest) (*User, *Response, error)
func (s *MemberServiceOp) Delete(UserID string) (*Response, error)
