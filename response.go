package icinga

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type apiResponse struct {
	Results []struct {
		Name   string
		Type   string
		Errors []string
		Attrs  json.RawMessage
	}
	Status string
}

type response struct {
	Results []object
	Error   error
}

func parseAPIResponse(r io.Reader) (apiResponse, error) {
	var apiresp apiResponse
	if err := json.NewDecoder(r).Decode(&apiresp); err != nil {
		return apiResponse{}, err
	}
	return apiresp, nil
}

func parseResponse(r io.Reader) (*response, error) {
	apiresp, err := parseAPIResponse(r)
	if err != nil {
		return nil, err
	}
	// Confusingly the top-level status field in an API response contains
	// an error message. Successful statuses are actually held in the
	// status field in Results!
	if apiresp.Status != "" {
		return &response{Error: errors.New(apiresp.Status)}, nil
	}
	resp := &response{}
	for _, r := range apiresp.Results {
		if len(r.Errors) > 0 {
			resp.Error = errors.New(strings.Join(r.Errors, ", "))
			// got an error so nothing left in the API response
			break
		}
		if r.Type == "" {
			continue //
		}
		switch r.Type {
		case "Host":
			var h Host
			if err := json.Unmarshal(r.Attrs, &h); err != nil {
				return nil, err
			}
			resp.Results = append(resp.Results, h)
		case "Service":
			var s Service
			if err := json.Unmarshal(r.Attrs, &s); err != nil {
				return nil, err
			}
			resp.Results = append(resp.Results, s)
		case "User":
			var u User
			if err := json.Unmarshal(r.Attrs, &u); err != nil {
				return nil, err
			}
			resp.Results = append(resp.Results, u)
		case "HostGroup":
			var h HostGroup
			if err := json.Unmarshal(r.Attrs, &h); err != nil {
				return nil, err
			}
			resp.Results = append(resp.Results, h)
		default:
			return nil, fmt.Errorf("unsupported unmarshal of type %s", r.Type)
		}
	}
	return resp, nil
}

func objectFromLookup(resp *response) (object, error) {
	if len(resp.Results) == 0 {
		return nil, errors.New("empty results")
	} else if len(resp.Results) > 1 {
		return nil, errors.New("too many results")
	}
	return resp.Results[0], nil
}
