package types

import (
	"fmt"
	"net/http"
)

// Response is the http response from api calls
type Response struct {
	*http.Response
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
