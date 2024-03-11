package latitude

const tagBasePath = "/tags"

type TagsService interface {
	List(*ListOptions) ([]Tag, *Response, error)
	Create(*TagCreateRequest) (*Tag, *Response, error)
	Update(string, *TagUpdateRequest) (*Tag, *Response, error)
	Delete(string) (*Response, error)
}

type Tag struct {
	ID          string
	Name        string
	Slug        string
	Description string
	Color       string
	TeamID      string
	TeamName    string
	TeamSlug    string
}

type TagListResponse struct {
	Data []TagData `json:"data"`
	Meta meta      `json:"meta"`
}

type TagResponse struct {
	Data TagData `json:"data"`
	Meta meta    `json:"meta"`
}

type TagData struct {
	ID         string        `json:"id"`
	Type       string        `json:"type"`
	Attributes TagAttributes `json:"name"`
}

type TagAttributes struct {
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	Color       string     `json:"color"`
	Team        ServerTeam `json:"team"`
}
type TagCreateRequest struct{}
type TagUpdateRequest struct{}

type TagServiceOp struct {
	client requestDoer
}

func NewFlatTag(td TagData) Tag {
	return Tag{
		ID:          td.ID,
		Name:        td.Attributes.Name,
		Slug:        td.Attributes.Slug,
		Description: td.Attributes.Description,
		Color:       td.Attributes.Color,
		TeamID:      td.Attributes.Team.ID,
		TeamName:    td.Attributes.Team.Name,
		TeamSlug:    td.Attributes.Team.Slug,
	}
}

func NewFlatTagList(td []TagData) []Tag {
	var res []Tag
	for _, tag := range td {
		res = append(res, NewFlatTag(tag))
	}
	return res
}

func (t *TagServiceOp) List(opts *ListOptions) ([]Tag, *Response, error) {
	apiPathQuery := opts.WithQuery(tagBasePath)
	var tags []Tag

	for {
		res := new(TagListResponse)

		resp, err := t.client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		tags = append(tags, NewFlatTagList(res.Data)...)

		if apiPathQuery = nextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}

		return tags, resp, err
	}
}

func (t *TagServiceOp) Create(createRequest *TagCreateRequest) (*Tag, *Response, error) {
	panic("")
}
func (t *TagServiceOp) Update(tagID string, updateRequest *TagUpdateRequest) (*Tag, *Response, error) {
	panic("")
}
func (t *TagServiceOp) Delete(tagID string) (*Response, error) {
	panic("")
}
