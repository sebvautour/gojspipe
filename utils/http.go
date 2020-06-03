package utils

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// HTTPRequest is given to HTTPReq
type HTTPRequest struct {
	Method   string
	URL      string
	Body     string
	Headers  map[string]string
	Username string
	Password string
}

// HTTPResponse is returned by HTTPReq
type HTTPResponse struct {
	StatusCode int
	Status     string
	Body       string
	// Headers    map[string]string
	Error string
}

// HTTPReq can be used to perform an HTTP request
func (u *Utils) HTTPReq(req HTTPRequest) (resp HTTPResponse) {

	var body io.Reader
	if req.Body != "" {
		body = strings.NewReader(req.Body)
	} else {
		body = nil
	}
	httpreq, err := http.NewRequest(req.Method, req.URL, body)
	if err != nil {
		return HTTPResponse{Error: err.Error()}
	}

	if req.Username != "" && req.Password != "" {
		httpreq.SetBasicAuth(req.Username, req.Password)
	}

	for k, v := range req.Headers {
		httpreq.Header.Set(k, v)
	}

	httpresp, err := u.httpClient.Do(httpreq)
	if err != nil {
		return HTTPResponse{Error: err.Error()}
	}

	b, err := ioutil.ReadAll(httpresp.Body)
	if err != nil {
		return HTTPResponse{Error: err.Error()}
	}
	defer httpresp.Body.Close()

	return HTTPResponse{
		StatusCode: httpresp.StatusCode,
		Status:     httpresp.Status,
		Body:       string(b),
		Error:      "",
	}
}
