package icinga

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type object interface {
	name() string
	path() string
}

// jsonForCreate marshals obj into the required JSON object to be sent
// in the body of a PUT request to Icinga. Some fields of obj must not be set for
// Icinga to create the object. Since some of those fields are structs
// (and not pointers to structs), they are always included, even if unset.
// jsonForCreate overrides those fields to always be empty. Other fields are left
// alone to let Icinga report an error for us.
func jsonForCreate(obj object) ([]byte, error) {
	m := make(map[string]interface{})
	switch v := obj.(type) {
	case User, HostGroup:
		m["attrs"] = v
	case Host:
		aux := &struct {
			// fields not added to Host yet
			// LastCheck struct{}
			// LastCheckResult struct{}
			Host
		}{Host: v}
		m["attrs"] = aux
	case Service:
		aux := &struct {
			LastCheck       *struct{} `json:",omitempty"`
			LastCheckResult *struct{} `json:"last_check_result,omitempty"`
			Service
		}{Service: v}
		m["attrs"] = aux
	default:
		return nil, fmt.Errorf("marshal %T for creation unsupported", v)
	}
	return json.Marshal(m)
}

//go:generate ./crud.sh -o crud.go

func (c *Client) lookupObject(objpath string) (object, error) {
	resp, err := c.get(objpath, "")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotExist
	}
	iresp, err := parseResponse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse response: %v", err)
	} else if iresp.Error != nil {
		return nil, iresp.Error
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	return objectFromLookup(iresp)
}

func (c *Client) filterObjects(objpath, expr string) ([]object, error) {
	resp, err := c.get(objpath, expr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	iresp, err := parseResponse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse response: %v", err)
	} else if iresp.Error != nil {
		return nil, iresp.Error
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	} else if len(iresp.Results) == 0 {
		return nil, ErrNoMatch
	}
	return iresp.Results, nil
}

func (c *Client) createObject(obj object) error {
	b, err := jsonForCreate(obj)
	if err != nil {
		return fmt.Errorf("marshal into json: %v", err)
	}
	resp, err := c.put(obj.path(), bytes.NewReader(b))
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	defer resp.Body.Close()
	iresp, err := parseResponse(resp.Body)
	if err != nil {
		return fmt.Errorf("parse response: %v", err)
	}
	if strings.Contains(iresp.Error.Error(), "already exists") {
		return ErrExist
	}
	return iresp.Error
}

func (c *Client) deleteObject(objpath string, cascade bool) error {
	resp, err := c.delete(objpath, cascade)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	} else if resp.StatusCode == http.StatusNotFound {
		return ErrNotExist
	}
	iresp, err := parseResponse(resp.Body)
	if err != nil {
		return fmt.Errorf("parse response: %v", err)
	}
	return iresp.Error
}
