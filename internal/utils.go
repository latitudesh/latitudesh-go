package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"

	types "github.com/latitudesh/latitudesh-go/types"
)

// meta contains pagination information
type Meta struct {
	Self           *Href `json:"self"`
	First          *Href `json:"first"`
	Last           *Href `json:"last"`
	Previous       *Href `json:"previous,omitempty"`
	Next           *Href `json:"next,omitempty"`
	Total          int   `json:"total"`
	CurrentPageNum int   `json:"current_page"`
	LastPageNum    int   `json:"last_page"`
}

// Href is an API link
type Href struct {
	Href string `json:"href"`
}

type RequestDoer interface {
	NewRequest(method, path string, body interface{}) (*http.Request, error)
	Do(req *http.Request, v interface{}) (*types.Response, error)
	DoRequest(method, path string, body, v interface{}) (*types.Response, error)
	DoRequestWithHeader(method string, headers map[string]string, path string, body, v interface{}) (*types.Response, error)
}

// from terraform-plugin-sdk/v2/helper/logging/transport.go
func PrettyPrintJsonLines(b []byte) string {
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

func DumpResponse(resp *http.Response) {
	o, _ := httputil.DumpResponse(resp, true)
	strResp := PrettyPrintJsonLines(o)
	reg, _ := regexp.Compile(`"token":(.+?),`)
	reMatches := reg.FindStringSubmatch(strResp)
	if len(reMatches) == 2 {
		strResp = strings.Replace(strResp, reMatches[1], strings.Repeat("-", len(reMatches[1])), 1)
	}
	log.Printf("\n=======[RESPONSE]============\n%s\n\n", strResp)
}

func DumpRequest(req *http.Request) {
	r := req.Clone(context.TODO())
	r.Body, _ = req.GetBody()
	h := r.Header
	if len(h.Get("Authorization")) != 0 {
		h.Set("Authorization", "**REDACTED**")
	}
	defer r.Body.Close()

	o, _ := httputil.DumpRequestOut(r, false)
	bbs, _ := io.ReadAll(r.Body)
	reqBodyStr := PrettyPrintJsonLines(bbs)
	strReq := PrettyPrintJsonLines(o)
	log.Printf("\n=======[REQUEST]=============\n%s%s\n", string(strReq), reqBodyStr)
}

// dumpDeprecation logs headers defined by
// https://tools.ietf.org/html/rfc8594
func DumpDeprecation(resp *http.Response) {
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

func CheckResponse(r *http.Response) error {

	if s := r.StatusCode; s >= 200 && s <= 299 {
		// response is good, return
		return nil
	}

	errorResponse := &types.ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	// if the response has a body, populate the message in errorResponse
	if err != nil {
		return err
	}

	if len(data) > 0 {
		err = json.Unmarshal(data, errorResponse)
		if err != nil {
			return err
		}
	}

	return errorResponse
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
