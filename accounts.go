package latitude

const userBasePath = "/user/profile"
const userTeamsPath = "/user/teams"

// UserService interface defines available Account methods
type UserService interface {
	Get(*GetOptions) (*User, *Response, error)
	Update(string, *UserUpdateRequest) (*User, *Response, error)
	List(listOpt *ListOptions) ([]Team, *Response, error)
}

// UserServiceOp implements UserService
type UserServiceOp struct {
	client requestDoer
}

type UserGetResponse struct {
	Data UserGetData `json:"data"`
	Meta meta        `json:"meta"`
}

type UserGetData struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes UserAttributes `json:"attributes"`
}

type UserAttributes struct {
	FirstName  string   `json:"first_name"`
	LastName   string   `json:"last_name"`
	Email      string   `json:"email"`
	MfaEnabled bool     `json:"mfa_enabled"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	Role       UserRole `json:"role"`
}

type UserUpdateRequest struct {
	Data UserUpdateData `json:"data"`
}

type UserUpdateData struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Attributes UserUpdateAttributes `json:"attributes"`
}

type UserUpdateAttributes struct {
	FirstName              string        `json:"first_name"`
	LastName               string        `json:"last_name"`
	Role                   AvailableRole `json:"role"`
	AuthenticationFactorId string        `json:"authentication_factor_id"`
}

type UserRole struct {
	Role
	createdAt string
	updatedAt string
}

type User struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	MfaEnabled bool   `json:"mfa_enabled"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Role       string `json:"role"`
}

// Flatten latitude API data structure
func NewFlatUser(md UserGetData) User {
	return User{
		md.ID,
		md.Attributes.FirstName,
		md.Attributes.LastName,
		md.Attributes.Email,
		md.Attributes.MfaEnabled,
		md.Attributes.CreatedAt,
		md.Attributes.UpdatedAt,
		md.Attributes.Role.Name,
	}
}

// Get the current User profile
func (s *UserServiceOp) Get(opts *GetOptions) (*User, *Response, error) {
	endpointPath := userBasePath
	apiPathQuery := opts.WithQuery(endpointPath)
	user := new(UserGetResponse)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, user)
	if err != nil {
		return nil, resp, err
	}

	flatUser := NewFlatUser(user.Data)
	return &flatUser, resp, err
}

// Update the User profile
func (s *UserServiceOp) Update(id string, updateRequest *UserUpdateRequest) (*User, *Response, error) {
	apiPath := userBasePath
	user := new(UserGetResponse)

	resp, err := s.client.DoRequest("PATCH", apiPath, updateRequest, user)
	if err != nil {
		return nil, resp, err
	}

	flatUser := NewFlatUser(user.Data)
	return &flatUser, resp, err
}

// List the current User teams
func (s *UserServiceOp) List(opts *ListOptions) ([]Team, *Response, error) {
	apiPathQuery := userTeamsPath
	teams := []Team{}

	for {
		res := new(TeamGetResponse)

		resp, err := s.client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		teams = append(teams, NewFlatTeamList(res.Data)...)

		if apiPathQuery = nextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}

		return teams, resp, nil
	}
}
