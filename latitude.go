package latitude

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	internal "github.com/latitudesh/latitudesh-go/internal"
	opsys "github.com/latitudesh/latitudesh-go/operating_systems"
	servers "github.com/latitudesh/latitudesh-go/servers"

	types "github.com/latitudesh/latitudesh-go/types"
)

const (
	authTokenEnvVar      = "LATITUDE_AUTH_TOKEN"
	apiVersion           = "2023-06-01"
	baseURL              = "https://api.latitude.sh"
	debugEnvVar          = "LATITUDE_DEBUG"
	userAgentForSDK      = "Latitude-Go-SDK"
	userAgentForProvider = "Latitude-Terraform-Provider"
)

var currentVersion = "0.3.1"

// Client is the base API Client
type Client struct {
	client        *http.Client
	debug         bool
	BaseURL       *url.URL
	UserAgent     string
	ConsumerToken string
	APIKey        string

	Projects         ProjectService
	Servers          servers.ServerService
	UserData         UserDataService
	SSHKeys          SSHKeyService
	Plans            PlanService
	OperatingSystems opsys.OperatingSystemService
	VirtualNetworks  VirtualNetworkService
	VlanAssignments  VlanAssignmentService
	Regions          RegionService
	Teams            TeamService
	Bandwidth        BandwidthService
	Members          MemberService
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

	req.Header.Add("API-Version", apiVersion)

	// set User-Agent value for SDK or terraform-provider
	userAgent := c.UserAgent
	if !strings.Contains(userAgent, userAgentForProvider) {
		userAgent = fmt.Sprintf("%s/%s", userAgentForSDK, currentVersion)
	}
	req.Header.Add("User-Agent", userAgent)

	if req.Method != "GET" {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

// Do executes the http request
func (c *Client) Do(req *http.Request, v interface{}) (*types.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := types.Response{Response: resp}

	if c.debug {
		internal.DumpResponse(response.Response)
	}
	internal.DumpDeprecation(response.Response)

	err = internal.CheckResponse(resp)
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

// DoRequest is a convenience method, it calls NewRequest followed by Do
// v is the interface to unmarshal the response JSON into
func (c *Client) DoRequest(method, path string, body, v interface{}) (*types.Response, error) {
	req, err := c.NewRequest(method, path, body)
	if c.debug {
		internal.DumpRequest(req)
	}
	if err != nil {
		return nil, err
	}
	return c.Do(req, v)
}

// DoRequestWithHeader same as DoRequest
func (c *Client) DoRequestWithHeader(method string, headers map[string]string, path string, body, v interface{}) (*types.Response, error) {
	req, err := c.NewRequest(method, path, body)
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if c.debug {
		internal.DumpRequest(req)
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
	c.Servers = &servers.ServerServiceOp{Client: c}
	c.SSHKeys = &SSHKeyServiceOp{client: c}
	c.UserData = &UserDataServiceOp{client: c}
	c.Teams = &TeamServiceOp{client: c}
	c.Bandwidth = &BandwidthServiceOp{client: c}
	c.Plans = &PlanServiceOp{client: c}
	c.OperatingSystems = &opsys.OperatingSystemServiceOp{Client: c}
	c.Regions = &RegionServiceOp{client: c}
	c.VirtualNetworks = &VirtualNetworkServiceOp{client: c}
	c.VlanAssignments = &VlanAssignmentServiceOp{client: c}
	c.Members = &MemberServiceOp{client: c}
	c.debug = os.Getenv(debugEnvVar) != ""

	return c, nil
}
