package latitude

import "path"

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

type EmbedTag struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
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
	Attributes TagAttributes `json:"attributes"`
}

type TagAttributes struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	Color       string  `json:"color"`
	Team        TagTeam `json:"team"`
}

type TagTeam struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Address     string   `json:"address"`
	Status      string   `json:"status"`
	Currency    Currency `json:"currency"`
}

type Currency struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type TagCreateRequest struct {
	Data TagCreateData `json:"data"`
}

type TagCreateData struct {
	Type       string              `json:"type"`
	Attributes TagCreateAttributes `json:"attributes"`
}

type TagCreateAttributes struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}

type TagUpdateRequest struct {
	Data TagUpdateData `json:"data"`
}

type TagUpdateData struct {
	ID         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes TagUpdateAttributes `json:"attributes"`
}

type TagUpdateAttributes TagCreateAttributes

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
	tag := new(TagResponse)

	resp, err := t.client.DoRequest("POST", tagBasePath, createRequest, tag)
	if err != nil {
		return nil, resp, err
	}

	flatTag := NewFlatTag(tag.Data)
	return &flatTag, resp, err
}

func (t *TagServiceOp) Update(tagID string, updateRequest *TagUpdateRequest) (*Tag, *Response, error) {
	apiPath := path.Join(tagBasePath, tagID)
	tag := new(TagResponse)

	resp, err := t.client.DoRequest("PATCH", apiPath, updateRequest, tag)
	if err != nil {
		return nil, resp, err
	}

	flatTag := NewFlatTag(tag.Data)
	return &flatTag, resp, err
}

func (t *TagServiceOp) Delete(tagID string) (*Response, error) {
	apiPath := path.Join(tagBasePath, tagID)

	return t.client.DoRequest("DELETE", apiPath, nil, nil)
}
