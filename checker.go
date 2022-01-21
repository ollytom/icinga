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

type checkFilter struct {
	Type string `json:"type"`
	Expr string `json:"filter"`
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

// Check reschedules the checks for all hosts in the HostGroup hg via the
// provided Client.
func (hg HostGroup) Check(c *Client) error {
	return c.check(hg)
}

func splitServiceName(name string) []string {
	return strings.SplitN(name, "!", 2)
}

func (c *Client) check(ch checker) error {
	switch v := ch.(type) {
	case Host:
		return c.CheckHosts(fmt.Sprintf("host.name == %q", v.Name))
	case Service:
		a := splitServiceName(v.Name)
		if len(a) != 2 {
			return fmt.Errorf("check %s: invalid service name", v.Name)
		}
		host := a[0]
		service := a[1]
		return c.CheckServices(fmt.Sprintf("host.name == %q && service.name == %q", host, service))
	case HostGroup:
		return c.CheckHosts(fmt.Sprintf("%q in host.groups", v.Name))
	default:
		return fmt.Errorf("cannot check %T", v)
	}
}

// CheckHosts schedules checks for all services matching the filter expression
// filter. If no services match the filter, error wraps ErrNoMatch.
func (c *Client) CheckServices(filter string) error {
	f := checkFilter{
		Type: "Service",
		Expr: filter,
	}
	if err := scheduleCheck(c, f); err != nil {
		return fmt.Errorf("check services %s: %w", filter, err)
	}
	return nil
}

// CheckHosts schedules checks for all hosts matching the filter expression
// filter. If no hosts match the filter, error wraps ErrNoMatch.
func (c *Client) CheckHosts(filter string) error {
	f := checkFilter{
		Type: "Host",
		Expr: filter,
	}
	if err := scheduleCheck(c, f); err != nil {
		return fmt.Errorf("check hosts %s: %w", filter, err)
	}
	return nil
}

func scheduleCheck(c *Client, filter checkFilter) error {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(filter); err != nil {
		return err
	}
	resp, err := c.post("/actions/reschedule-check", buf)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	} else if resp.StatusCode == http.StatusNotFound {
		return ErrNoMatch
	}
	defer resp.Body.Close()
	iresp, err := parseResponse(resp.Body)
	if err != nil {
		return fmt.Errorf("parse response: %v", err)
	}
	if iresp.Error != nil {
		return iresp.Error
	}
	return fmt.Errorf("%s", resp.Status)
}
