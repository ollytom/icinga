package icinga

import (
	"encoding/json"
	"errors"
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
	Attrs  map[string]interface{}
	Code   int
	Errors []string
	Name   string
	Type   string
}

var ErrNoObject = errors.New("no such object")

func (res results) Err() error {
	if len(res.Results) == 0 {
		return nil
	}
	var errs []string
	for _, r := range res.Results {
		if len(r.Errors) == 0 {
			continue
		}
		errs = append(errs, strings.Join(r.Errors, ", "))
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.New(strings.Join(errs, ", "))
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
	req, err := newRequest(http.MethodGet, c.addr, path, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) post(path string, body io.Reader) (*http.Response, error) {
	req, err := newRequest(http.MethodPost, c.addr, path, body)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) put(path string, body io.Reader) error {
	req, err := newRequest(http.MethodPut, c.addr, path, body)
	if err != nil {
		return err
	}
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	defer resp.Body.Close()
	var results results
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return results.Err()
}

func (c *Client) delete(path string) error {
	req, err := newRequest(http.MethodDelete, c.addr, path, nil)
	if err != nil {
		return err
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
	return results.Err()
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(c.username, c.password)
	return c.Do(req)
}
