package icinga

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func (c *Client) get(path string) (*http.Response, error) {
	url := "https://" + c.addr + versionPrefix + path
	req, err := NewRequest(http.MethodGet, url, c.username, c.password, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) getFilter(path, filter string) (*http.Response, error) {
	u, err := url.Parse("https://" + c.addr + versionPrefix + path)
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("filter", filter)
	u.RawQuery = v.Encode()
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

func (c *Client) delete(path string) (*http.Response, error) {
	url := "https://" + c.addr + versionPrefix + path
	req, err := NewRequest(http.MethodDelete, url, c.username, c.password, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
