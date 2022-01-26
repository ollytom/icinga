package icinga

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const versionPrefix = "/v1"

// NewRequest returns an authenticated HTTP request with appropriate header
// for sending to an Icinga2 server.
func NewRequest(method, url, username, password string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	switch req.Method {
	case http.MethodGet, http.MethodDelete:
		req.Header.Set("Accept", "application/json")
	case http.MethodPost, http.MethodPut:
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
	default:
		return nil, fmt.Errorf("new request: unsupported method %s", req.Method)
	}
	req.SetBasicAuth(username, password)
	return req, nil
}

// filterEncode url-encodes the filter expression expr in the required
// format to be included in a request. Notably, spaces need to be encoded
// as "%20", not "+" as returned by the url package.
func filterEncode(expr string) string {
	v := url.Values{}
	v.Set("filter", expr)
	return strings.ReplaceAll(v.Encode(), "+", "%20")
}

func (c *Client) get(path, filter string) (*http.Response, error) {
	u, err := url.Parse("https://" + c.addr + versionPrefix + path)
	if err != nil {
		return nil, err
	}
	if filter != "" {
		u.RawQuery = filterEncode(filter)
	}
	req, err := NewRequest(http.MethodGet, u.String(), c.username, c.password, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) post(path string, body io.Reader) (*http.Response, error) {
	url := "https://" + c.addr + versionPrefix + path
	req, err := NewRequest(http.MethodPost, url, c.username, c.password, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) put(path string, body io.Reader) (*http.Response, error) {
	url := "https://" + c.addr + versionPrefix + path
	req, err := NewRequest(http.MethodPut, url, c.username, c.password, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) delete(path string, cascade bool) (*http.Response, error) {
	u, err := url.Parse("https://" + c.addr + versionPrefix + path)
	if err != nil {
		return nil, err
	}
	if cascade {
		v := url.Values{}
		v.Set("cascade", "1")
		u.RawQuery = v.Encode()
	}
	req, err := NewRequest(http.MethodDelete, u.String(), c.username, c.password, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
