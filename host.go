package icinga

import "encoding/json"

// Host represents a Host object. To create a Host, the Name and CheckCommand
// fields must be set.
type Host struct {
	Name            string    `json:"-"`
	Address         string    `json:"address"`
	Address6        string    `json:"address6"`
	Groups          []string  `json:"groups,omitempty"`
	State           HostState `json:"state,omitempty"`
	StateType       StateType `json:"state_type,omitempty"`
	CheckCommand    string    `json:"check_command"`
	DisplayName     string    `json:"display_name,omitempty"`
	Acknowledgement bool      `json:",omitempty"`
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
		Acknowledgement int
		*alias
	}{
		alias: (*alias)(h),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Acknowledgement != 0 {
		h.Acknowledgement = true
	}
	return nil
}
