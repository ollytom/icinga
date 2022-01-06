package icinga

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const versionPrefix = "/v1"

type results struct {
	Results []result
}

type result struct {
	Attrs  interface{}
	Code   int
	Errors []string
	Name   string
	Type   string
}

var ErrNoObject = errors.New("no such object")

func (res results) Error() string {
	var s []string
	for _, r := range res.Results {
		s = append(s, r.Error())
	}
	return strings.Join(s, ", ")
}

func (r result) Error() string {
	return strings.Join(r.Errors, ", ")
}

func newRequest(method, host, path string, body io.Reader) (*http.Request, error) {
	url := "https://" + host + versionPrefix + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	switch req.Method {
	case http.MethodGet:
		break
	case http.MethodDelete:
		req.Header.Set("Accept", "application/json")
	case http.MethodPost, http.MethodPut:
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
	default:
		return nil, fmt.Errorf("new request: unsupported method %s", req.Method)
	}
	return req, nil
}

func (c *Client) get(path string) (*http.Response, error) {
	req, err := newRequest(http.MethodGet, c.host, path, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) post(path string, body io.Reader) (*http.Response, error) {
	req, err := newRequest(http.MethodPost, c.host, path, body)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) put(path string, body io.Reader) (*http.Response, error) {
	req, err := newRequest(http.MethodPut, c.host, path, body)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) delete(path string) error {
	req, err := newRequest(http.MethodDelete, c.host, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	} else if resp.StatusCode == http.StatusNotFound {
		return ErrNoObject
	}
	defer resp.Body.Close()
	var results results
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return results
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(c.username, c.password)
	return c.Do(req)
}
