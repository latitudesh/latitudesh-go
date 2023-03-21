package latitude

import "path"

const teamBasePath = "/team"

type TeamService interface {
	Get() (*Team, *Response, error)
	Create(request *TeamCreateRequest) (*Team, *Response, error)
	Update(TeamID string, request *TeamUpdateRequest) (*Team, *Response, error)
}

// Team represents a Latitude Team record
type Team struct {
	ID          string        `json:"id"`
	Description string        `json:"description"`
	Name        string        `json:"name"`
	Slug        string        `json:"slug"`
	Currency    string        `json:"currency"`
	Address     *string       `json:"address"`
	Status      *string       `json:"status"`
	Projects    []interface{} `json:"projects"`
	Users       []interface{} `json:"users"`
	Owner       interface{}   `json:"owner"`
	Billing     interface{}   `json:"billing"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
}

// TeamServiceOp implements TeamService
type TeamServiceOp struct {
	client requestDoer
}

// TeamCreateRequest type used to create a Latitude Team record
type TeamCreateRequest struct {
	Data TeamCreateData `json:"data"`
}

type TeamCreateData struct {
	Type       string               `json:"type"`
	Attributes TeamCreateAttributes `json:"attributes"`
}

type TeamCreateAttributes struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	Currency    string `json:"currency"`
	Address     string `json:"address,omitempty"`
}

type TeamUpdateRequest struct {
	Data TeamUpdateData `json:"data"`
}

type TeamUpdateData struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Attributes TeamUpdateAttributes `json:"attributes"`
}

type TeamUpdateAttributes struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	Address     string `json:"address,omitempty"`
}

type TeamGetResponse struct {
	Data []TeamData `json:"data"`
	Meta meta       `json:"meta"`
}

type TeamCreateResponse struct {
	Data TeamData `json:"data"`
	Meta meta     `json:"meta"`
}

type TeamData struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes TeamGetAttributes `json:"attributes"`
}

type TeamGetAttributes struct {
	Description string        `json:"description"`
	Name        string        `json:"name"`
	Slug        string        `json:"slug"`
	Currency    string        `json:"currency"`
	Address     *string       `json:"address"`
	Status      *string       `json:"status"`
	Projects    []interface{} `json:"projects"`
	Users       []interface{} `json:"users"`
	Owner       interface{}   `json:"owner"`
	Billing     interface{}   `json:"billing"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
}

// Flatten latitude API data structures
func NewFlatTeam(t TeamData) Team {
	return Team{
		t.ID,
		t.Attributes.Description,
		t.Attributes.Name,
		t.Attributes.Slug,
		t.Attributes.Currency,
		t.Attributes.Address,
		t.Attributes.Status,
		t.Attributes.Projects,
		t.Attributes.Users,
		t.Attributes.Owner,
		t.Attributes.Billing,
		t.Attributes.CreatedAt,
		t.Attributes.UpdatedAt,
	}
}

// Get returns a Team by id
func (u *TeamServiceOp) Get() (*Team, *Response, error) {
	var flatTeam Team
	Team := new(TeamGetResponse)

	resp, err := u.client.DoRequest("GET", teamBasePath, nil, Team)
	if err != nil {
		return nil, resp, err
	}

	for _, team := range Team.Data {
		flatTeam = NewFlatTeam(team)
	}
	return &flatTeam, resp, err
}

// Create creates a new Team record
func (s *TeamServiceOp) Create(createRequest *TeamCreateRequest) (*Team, *Response, error) {
	Team := new(TeamCreateResponse)

	resp, err := s.client.DoRequest("POST", teamBasePath, createRequest, Team)
	if err != nil {
		return nil, resp, err
	}

	flatTeam := NewFlatTeam(Team.Data)
	return &flatTeam, resp, err
}

// Update updates a Team record
func (s *TeamServiceOp) Update(TeamID string, updateRequest *TeamUpdateRequest) (*Team, *Response, error) {
	apiPath := path.Join(teamBasePath, TeamID)
	Team := new(TeamCreateResponse)

	resp, err := s.client.DoRequest("PATCH", apiPath, updateRequest, Team)
	if err != nil {
		return nil, resp, err
	}

	flatTeam := NewFlatTeam(Team.Data)
	return &flatTeam, resp, err
}
