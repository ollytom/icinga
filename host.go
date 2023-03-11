package icinga

import (
	"encoding/json"
	"fmt"
	"time"
)

// Host represents a Host object. To create a Host, the Name and CheckCommand
// fields must be set.
type Host struct {
	Name            string      `json:"-"`
	Address         string      `json:"address"`
	Address6        string      `json:"address6"`
	Groups          []string    `json:"groups,omitempty"`
	State           HostState   `json:"state,omitempty"`
	StateType       StateType   `json:"state_type,omitempty"`
	CheckCommand    string      `json:"check_command"`
	DisplayName     string      `json:"display_name,omitempty"`
	LastCheck       time.Time   `json:",omitempty"`
	LastCheckResult CheckResult `json:"last_check_result,omitempty"`
	Acknowledgement bool        `json:",omitempty"`
	Notes           string      `json:"notes,omitempty"`
	NotesURL        string      `json:"notes_url,omitempty"`
}

type HostGroup struct {
	Name        string `json:"-"`
	DisplayName string `json:"display_name"`
}

type HostState int

const (
	HostUp HostState = 0 + iota
	HostDown
	HostUnreachable
)

func (state HostState) String() string {
	switch state {
	case HostUp:
		return "HostUp"
	case HostDown:
		return "HostDown"
	}
	return "HostUnreachable"
}

func (h Host) name() string {
	return h.Name
}

func (h Host) path() string {
	return "/objects/hosts/" + h.Name
}

func (hg HostGroup) name() string {
	return hg.Name
}

func (hg HostGroup) path() string {
	return "/objects/hostgroups/" + hg.Name
}

// UnmarhsalJSON unmarshals host attributes into more meaningful Host field types.
func (h *Host) UnmarshalJSON(data []byte) error {
	type alias Host
	aux := &struct {
		Acknowledgement interface{} `json:"acknowledgement"`
		State           interface{} `json:"state"`
		StateType       interface{} `json:"state_type"`
		LastCheck       float64     `json:"last_check"`
		*alias
	}{
		alias: (*alias)(h),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		fmt.Println("uh oh!")
		return err
	}
	switch v := aux.Acknowledgement.(type) {
	case int:
		if v != 0 {
			h.Acknowledgement = true
		}
	case float64:
		if int(v) != 0 {
			h.Acknowledgement = true
		}
	}
	switch v := aux.State.(type) {
	case int:
		h.State = HostState(v)
	case float64:
		h.State = HostState(v)
	}
	switch v := aux.StateType.(type) {
	case int:
		h.StateType = StateType(v)
	case float64:
		h.StateType = StateType(v)
	}
	h.LastCheck = time.Unix(int64(aux.LastCheck), 0)
	return nil
}
