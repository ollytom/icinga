package icinga

import "encoding/json"

// Host represents a Host object. To create a Host, the Name and CheckCommand
// fields must be set.
type Host struct {
	Name            string    `json:"name"`
	Address         string    `json:"address"`
	Address6        string    `json:"address6"`
	Groups          []string  `json:"groups"`
	State           HostState `json:"state"`
	StateType       StateType `json:"state_type"`
	CheckCommand    string    `json:"check_command"`
	DisplayName     string    `json:"display_name"`
	Acknowledgement bool
}

type HostGroup struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
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

func (hg HostGroup) name() string {
	return hg.Name
}

func (hg HostGroup) path() string {
	return "/objects/hostgroups/" + hg.Name
}

func (h Host) MarshalJSON() ([]byte, error) {
	attrs := make(map[string]interface{})
	attrs["address"] = h.Address
	attrs["address6"] = h.Address6
	if len(h.Groups) > 0 {
		attrs["groups"] = h.Groups
	}
	attrs["check_command"] = h.CheckCommand
	attrs["display_name"] = h.DisplayName
	m := make(map[string]interface{})
	m["attrs"] = attrs
	return json.Marshal(m)
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

func (hg HostGroup) MarshalJSON() ([]byte, error) {
	type attrs struct {
		DisplayName string `json:"display_name"`
	}
	type group struct {
		Attrs attrs `json:"attrs"`
	}
	return json.Marshal(&group{
		Attrs: attrs{
			DisplayName: hg.DisplayName,
		},
	})
}
