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

// NewRequest returns an authenticated HTTP request with appropriate header
// for sending to an Icinga2 server.
func NewRequest(method, url, username, password string, body io.Reader) (*http.Request, error) {
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
	req.SetBasicAuth(username, password)
	return req, nil
}

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

func (c *Client) get(path string) (*http.Response, error) {
	url := "https://" + c.addr + versionPrefix + path
	req, err := NewRequest(http.MethodGet, url, c.username, c.password, nil)
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

func (c *Client) put(path string, body io.Reader) error {
	url := "https://" + c.addr + versionPrefix + path
	req, err := NewRequest(http.MethodPost, url, c.username, c.password, body)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
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
	url := "https://" + c.addr + versionPrefix + path
	req, err := NewRequest(http.MethodPost, url, c.username, c.password, body)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
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
