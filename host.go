package icinga

import "encoding/json"

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
