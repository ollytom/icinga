package icinga

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type object interface {
	name() string
	path() string
}

//go:generate ./crud.sh -o crud.go

func (c *Client) lookupObject(objpath string) (object, error) {
	return lookup(c, objpath)
}

func (c *Client) filterObjects(objpath, expr string) ([]object, error) {
	return filter(c, objpath, expr)
}

func (c *Client) createObject(obj object) error {
	return create(c, obj)
}

func (c *Client) deleteObject(objpath string, cascade bool) error {
	return delete(c, objpath, cascade)
}

type client interface {
	get(path, filter string) (*http.Response, error)
	post(path string, body io.Reader) (*http.Response, error)
	put(path string, body io.Reader) (*http.Response, error)
	delete(path string, cascade bool) (*http.Response, error)
}

func lookup(c client, objpath string) (object, error) {
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

func filter(c client, objpath, expr string) ([]object, error) {
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

func create(c client, obj object) error {
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

func delete(c client, objpath string, cascade bool) error {
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
