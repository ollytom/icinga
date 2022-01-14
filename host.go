package icinga

import "encoding/json"

// Host represents a Host object. To create a Host, the Name and CheckCommand
// fields must be set.
type Host struct {
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	Address6     string    `json:"address6"`
	Groups       []string  `json:"groups"`
	State        HostState `json:"state"`
	CheckCommand string    `json:"check_command"`
	DisplayName  string    `json:"display_name"`
}

type HostGroup struct {
	Name string `json:"name"`
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

func (hg HostGroup) MarshalJSON() ([]byte, error) {
	type attrs struct {
		DisplayName  string `json:"display_name"`
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
