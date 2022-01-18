package icinga

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type checker interface {
	object
	Check(*Client) error
}

type StateType int

const (
	StateSoft StateType = 0 + iota
	StateHard
)

func (st StateType) String() string {
	switch st {
	case StateSoft:
		return "StateSoft"
	case StateHard:
		return "StateHard"
	}
	return "unsupported state type"
}

// Check reschedules the check for s via the provided Client.
func (s Service) Check(c *Client) error {
	return c.check(s)
}

// Check reschedules the check for h via the provided Client.
func (h Host) Check(c *Client) error {
	return c.check(h)
}

func splitServiceName(name string) []string {
	return strings.SplitN(name, "!", 2)
}

func (c *Client) check(ch checker) error {
	var filter struct {
		Type string `json:"type"`
		Expr string `json:"filter"`
	}
	switch v := ch.(type) {
	case Host:
		filter.Type = "Host"
		filter.Expr = fmt.Sprintf("host.name == %q", v.Name)
	case Service:
		filter.Type = "Service"
		a := splitServiceName(v.Name)
		if len(a) != 2 {
			return fmt.Errorf("check %s: invalid service name", v.Name)
		}
		host := a[0]
		service := a[1]
		filter.Expr = fmt.Sprintf("host.name == %q && service.name == %q", host, service)
	default:
		return fmt.Errorf("cannot check %T", v)
	}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(filter); err != nil {
		return err
	}
	resp, err := c.post("/actions/reschedule-check", buf)
	if err != nil {
		return fmt.Errorf("check %s: %w", ch.name(), err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("check %s: %w", ch.name(), ErrNotExist)
	default:
		return fmt.Errorf("check %s: %s", ch.name(), resp.Status)
	}
}
