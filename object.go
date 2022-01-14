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
	if expr != "" && resp.StatusCode == http.StatusNotFound {
		return nil, ErrNoMatch

	}
	iresp, err := parseResponse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse response: %v", err)
	} else if iresp.Error != nil {
		return nil, iresp.Error
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	return iresp.Results, nil
}

func (c *Client) createObject(obj object) error {
	buf := &bytes.Buffer{}
	switch v := obj.(type) {
	case Host, Service, User, HostGroup:
		if err := json.NewEncoder(buf).Encode(v); err != nil {
			return err
		}
	default:
		return fmt.Errorf("create type %T unsupported", v)
	}
	resp, err := c.put(obj.path(), buf)
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

func (c *Client) deleteObject(objpath string) error {
	resp, err := c.delete(objpath)
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
