package icinga

import (
	"encoding/json"
	"io"
	"net/http"
)

const versionPrefix = "/v1"

func newRequest(method, host, path string, body io.Reader) (*http.Request, error) {
	url := "https://" + host + versionPrefix + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	switch req.Method {
	case http.MethodPost, http.MethodPut:
		req.Header.Set("Content-Type", "application/json")
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

func (c *Client) delete(path string, body io.Reader) (*http.Response, error) {
	req, err := newRequest(http.MethodDelete, c.host, path, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(c.username, c.password)
	return c.httpClient.Do(req)
}

// TODO
type Result struct{}
func results(resp *http.Response) (Result, error) {
	defer resp.Body.Close()
	var results []Result
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return Result{}, err
	}
	return Result{}, nil
}
