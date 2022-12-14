package latitude

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const (
	authTokenEnvVar = "LATITUDE_AUTH_TOKEN"
	baseURL         = "https://api.maxihost.com"
	debugEnvVar     = "LATITUDE_DEBUG"
)

// meta contains pagination information
type meta struct {
	Self           *Href `json:"self"`
	First          *Href `json:"first"`
	Last           *Href `json:"last"`
	Previous       *Href `json:"previous,omitempty"`
	Next           *Href `json:"next,omitempty"`
	Total          int   `json:"total"`
	CurrentPageNum int   `json:"current_page"`
	LastPageNum    int   `json:"last_page"`
}

// Response is the http response from api calls
type Response struct {
	*http.Response
}

// Href is an API link
type Href struct {
	Href string `json:"href"`
}

// ErrorResponse is the http response used on errors
type ErrorResponse struct {
	Response *http.Response
	Errors   []ErrorData `json:"errors,omitempty"`
}

type ErrorData struct {
	Code   string `json:"code"`
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (r *ErrorResponse) Error() string {
	err := ""
	for _, e := range r.Errors {
		err += fmt.Sprintf("%v %v: %d\n\n%v\nCODE: %v\nSTATUS: %v\nDETAIL: %v\n",
			r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, e.Title, e.Code, e.Status, e.Detail)
	}
	return err
}

// Client is the base API Client
type Client struct {
	client        *http.Client
	debug         bool
	BaseURL       *url.URL
	UserAgent     string
	ConsumerToken string
	APIKey        string

	Projects ProjectService
	Servers  ServerService
	SSHKeys  SSHKeyService
	UserData UserDataService
	Plans    PlanService
	Regions  RegionService
}

type requestDoer interface {
	NewRequest(method, path string, body interface{}) (*http.Request, error)
	Do(req *http.Request, v interface{}) (*Response, error)
	DoRequest(method, path string, body, v interface{}) (*Response, error)
	DoRequestWithHeader(method string, headers map[string]string, path string, body, v interface{}) (*Response, error)
}

// NewRequest inits a new http request with the proper headers
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	// relative path to append to the endpoint url, no leading slash please
	if path[0] == '/' {
		path = path[1:]
	}
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	// json encode the request body, if any
	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Close = true

	req.Header.Add("Authorization", c.APIKey)
	req.Header.Add("User-Agent", c.UserAgent)
	if req.Method != "GET" {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

// Do executes the http request
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := Response{Response: resp}

	if c.debug {
		dumpResponse(response.Response)
	}
	dumpDeprecation(response.Response)

	err = checkResponse(resp)
	// if the response is an error, return the ErrorResponse
	if err != nil {
		return &response, err
	}

	if v != nil {
		// if v implements the io.Writer interface, return the raw response
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return &response, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return &response, err
			}
		}
	}

	return &response, err
}

// dumpDeprecation logs headers defined by
// https://tools.ietf.org/html/rfc8594
func dumpDeprecation(resp *http.Response) {
	uri := ""
	if resp.Request != nil {
		uri = resp.Request.Method + " " + resp.Request.URL.Path
	}

	deprecation := resp.Header.Get("Deprecation")
	if deprecation != "" {
		if deprecation == "true" {
			deprecation = ""
		} else {
			deprecation = " on " + deprecation
		}
		log.Printf("WARNING: %q reported deprecation%s", uri, deprecation)
	}

	sunset := resp.Header.Get("Sunset")
	if sunset != "" {
		log.Printf("WARNING: %q reported sunsetting on %s", uri, sunset)
	}

	links := resp.Header.Values("Link")

	for _, s := range links {
		for _, ss := range strings.Split(s, ",") {
			if strings.Contains(ss, "rel=\"sunset\"") {
				link := strings.Split(ss, ";")[0]
				log.Printf("WARNING: See %s for sunset details", link)
			} else if strings.Contains(ss, "rel=\"deprecation\"") {
				link := strings.Split(ss, ";")[0]
				log.Printf("WARNING: See %s for deprecation details", link)
			}
		}
	}
}

// from terraform-plugin-sdk/v2/helper/logging/transport.go
func prettyPrintJsonLines(b []byte) string {
	parts := strings.Split(string(b), "\n")
	for i, p := range parts {
		if b := []byte(p); json.Valid(b) {
			var out bytes.Buffer
			_ = json.Indent(&out, b, "", " ")
			parts[i] = out.String()
		}
	}
	return strings.Join(parts, "\n")
}

func dumpResponse(resp *http.Response) {
	o, _ := httputil.DumpResponse(resp, true)
	strResp := prettyPrintJsonLines(o)
	reg, _ := regexp.Compile(`"token":(.+?),`)
	reMatches := reg.FindStringSubmatch(strResp)
	if len(reMatches) == 2 {
		strResp = strings.Replace(strResp, reMatches[1], strings.Repeat("-", len(reMatches[1])), 1)
	}
	log.Printf("\n=======[RESPONSE]============\n%s\n\n", strResp)
}

func dumpRequest(req *http.Request) {
	r := req.Clone(context.TODO())
	r.Body, _ = req.GetBody()
	h := r.Header
	if len(h.Get("Authorization")) != 0 {
		h.Set("Authorization", "**REDACTED**")
	}
	defer r.Body.Close()

	o, _ := httputil.DumpRequestOut(r, false)
	bbs, _ := ioutil.ReadAll(r.Body)
	reqBodyStr := prettyPrintJsonLines(bbs)
	strReq := prettyPrintJsonLines(o)
	log.Printf("\n=======[REQUEST]=============\n%s%s\n", string(strReq), reqBodyStr)
}

// DoRequest is a convenience method, it calls NewRequest followed by Do
// v is the interface to unmarshal the response JSON into
func (c *Client) DoRequest(method, path string, body, v interface{}) (*Response, error) {
	req, err := c.NewRequest(method, path, body)
	if c.debug {
		dumpRequest(req)
	}
	if err != nil {
		return nil, err
	}
	return c.Do(req, v)
}

// DoRequestWithHeader same as DoRequest
func (c *Client) DoRequestWithHeader(method string, headers map[string]string, path string, body, v interface{}) (*Response, error) {
	req, err := c.NewRequest(method, path, body)
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if c.debug {
		dumpRequest(req)
	}
	if err != nil {
		return nil, err
	}
	return c.Do(req, v)
}

// NewClient initializes and returns a Client
func NewClient() (*Client, error) {
	apiToken := os.Getenv(authTokenEnvVar)
	if apiToken == "" {
		return nil, fmt.Errorf("you must export %s", authTokenEnvVar)
	}
	c := NewClientWithAuth("latitude lib", apiToken, nil)
	return c, nil

}

// NewClientWithAuth initializes and returns a Client
func NewClientWithAuth(consumerToken string, apiKey string, httpClient *http.Client) *Client {
	client, _ := NewClientWithBaseURL(apiKey, httpClient, baseURL)
	return client
}

// NewClientWithBaseURL returns a Client pointing to nonstandard API URL, e.g.
// for mocking the remote API
func NewClientWithBaseURL(apiKey string, httpClient *http.Client, apiBaseURL string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	u, err := url.Parse(apiBaseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{client: httpClient, BaseURL: u, APIKey: apiKey}
	c.Projects = &ProjectServiceOp{client: c}
	c.Servers = &ServerServiceOp{client: c}
	c.SSHKeys = &SSHKeyServiceOp{client: c}
	c.UserData = &UserDataServiceOp{client: c}
	c.Plans = &PlanServiceOp{client: c}
	c.Regions = &RegionServiceOp{client: c}
	c.debug = os.Getenv(debugEnvVar) != ""

	return c, nil
}

func checkResponse(r *http.Response) error {

	if s := r.StatusCode; s >= 200 && s <= 299 {
		// response is good, return
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	// if the response has a body, populate the message in errorResponse
	if err != nil {
		return err
	}

	// ct := r.Header.Get("Content-Type")
	// if !strings.HasPrefix(ct, expectedAPIContentTypePrefix) {
	// 	errorResponse.SingleError = fmt.Sprintf("Unexpected Content-Type %s with status %s", ct, r.Status)
	// 	return errorResponse
	// }

	if len(data) > 0 {
		err = json.Unmarshal(data, errorResponse)
		if err != nil {
			return err
		}
	}

	return errorResponse
}
