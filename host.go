package icinga

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Host represents a Host object.
type Host struct {
	Name         string   `json:"name"`
	Address      string   `json:"address"`
	Address6     string   `json:"address6"`
	Groups       []string `json:"groups"`
	State        int      `json:"state"`
	CheckCommand string   `json:"check_command"`
	DisplayName  string   `json:"display_name"`
}

type hostresults struct {
	Results []hostresult `json:"results"`
	results
}

type hostresult struct {
	Host Host `json:"attrs"`
	result
}

var ErrNoHost = errors.New("no such host")

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

// Hosts returns all Hosts in the Icinga2 configuration.
func (c *Client) Hosts() ([]Host, error) {
	resp, err := c.get("/objects/hosts")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var res hostresults
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	var hosts []Host
	for _, r := range res.Results {
		hosts = append(hosts, r.Host)
	}
	return hosts, nil
}

// LookupHost returns the Host identified by name.
// If no Host is found, error wraps ErrNoHost.
func (c *Client) LookupHost(name string) (Host, error) {
	resp, err := c.get("/objects/hosts/" + name)
	if err != nil {
		return Host{}, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return Host{}, fmt.Errorf("lookup %s: %w", name, ErrNoHost)
	}
	return Host{}, err
}

// CreateHost creates the Host host.
// The Name and CheckCommand fields of host must be non-zero.
func (c *Client) CreateHost(host Host) error {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(host); err != nil {
		return err
	}
	if err := c.put("/objects/hosts/"+host.Name, buf); err != nil {
		return fmt.Errorf("create host %s: %w", host.Name, err)
	}
	return nil
}

// DeleteHost deletes the Host identified by name.
// If no Host is found, error wraps ErrNoObject.
func (c *Client) DeleteHost(name string) error {
	if err := c.delete("/objects/hosts/" + name); err != nil {
		return fmt.Errorf("delete host %s: %w", name, err)
	}
	return nil
}
