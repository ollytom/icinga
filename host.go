package icinga

import (
	"encoding/json"
	"fmt"
)

// Host represents a Host object.
type Host struct {
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	Address6     string    `json:"address6"`
	Groups       []string  `json:"groups"`
	State        HostState `json:"state"`
	CheckCommand string    `json:"check_command"`
	DisplayName  string    `json:"display_name"`
}

type HostState int

const (
	HostUp HostState = 0 + iota
	HostDown
	HostUnreachable
)

func (s HostState) String() string {
	switch s {
	case HostUp:
		return "HostUp"
	case HostDown:
		return "HostDown"
	case HostUnreachable:
		return "HostUnreachable"
	}
	return "unhandled host state"
}

func (h Host) name() string {
	return h.Name
}

func (h Host) path() string {
	return "/objects/hosts/" + h.Name
}

func (h Host) attrs() map[string]interface{} {
	m := make(map[string]interface{})
	m["display_name"] = h.DisplayName
	return m
}

func (h Host) MarshalJSON() ([]byte, error) {
	type Attrs struct {
		Address      string `json:"address"`
		CheckCommand string `json:"check_command"`
		DisplayName  string `json:"display_name"`
	}
	type host struct {
		Attrs Attrs `json:"attrs"`
	}
	jhost := &host{
		Attrs: Attrs{
			Address:      h.Address,
			CheckCommand: h.CheckCommand,
			DisplayName:  h.DisplayName,
		},
	}
	return json.Marshal(jhost)
}

// Hosts returns Hosts matching the filter expression filter.
// If no hosts match, error wraps ErrNoMatch.
// To fetch all hosts, set filter to the empty string ("").
func (c *Client) Hosts(filter string) ([]Host, error) {
	objects, err := c.filterObjects("/objects/hosts", filter)
	if err != nil {
		return nil, fmt.Errorf("get hosts filter %q: %w", filter, err)
	}
	var hosts []Host
	for _, o := range objects {
		v, ok := o.(Host)
		if !ok {
			return nil, fmt.Errorf("get all hosts: %T in response", v)
		}
		hosts = append(hosts, v)
	}
	return hosts, nil
}

// LookupHost returns the Host identified by name. If no Host is found,
// error wraps ErrNotExist.
func (c *Client) LookupHost(name string) (Host, error) {
	obj, err := c.lookupObject("/objects/hosts/" + name)
	if err != nil {
		return Host{}, fmt.Errorf("lookup %s: %w", name, err)
	}
	v, ok := obj.(Host)
	if !ok {
		return Host{}, fmt.Errorf("lookup %s: result type %T is not host", name, obj)
	}
	return v, nil
}

// CreateHost creates the Host host.
// The Name and CheckCommand fields of host must be non-zero.
func (c *Client) CreateHost(host Host) error {
	if err := c.createObject(host); err != nil {
		return fmt.Errorf("create host %s: %w", host.Name, err)
	}
	return nil
}

// DeleteHost deletes the Host identified by name.
// If no Host is found, error wraps ErrNotExist.
func (c *Client) DeleteHost(name string) error {
	if err := c.deleteObject("/objects/hosts/" + name); err != nil {
		return fmt.Errorf("delete host %s: %w", name, err)
	}
	return nil
}
