package icinga

import (
	"encoding/json"
	"strings"
	"time"
)

func (s Service) name() string {
	return s.Name
}

func (s Service) path() string {
	return "/objects/services/" + s.Name
}

// Service represents a Service object.
type Service struct {
	Name            string       `json:"-"`
	Groups          []string     `json:"groups,omitempty"`
	State           ServiceState `json:"state,omitempty"`
	StateType       StateType    `json:"state_type,omitempty"`
	CheckCommand    string       `json:"check_command"`
	DisplayName     string       `json:"display_name,omitempty"`
	LastCheck       time.Time    `json:",omitempty"`
	LastCheckResult *CheckResult `json:"last_check_result,omitempty"`
	Acknowledgement bool         `json:",omitempty"`
}

type CheckResult struct {
	CheckSource string `json:"check_source"`
	Command     interface{}
	Output      string
}

type ServiceState int

const (
	ServiceOK ServiceState = 0 + iota
	ServiceWarning
	ServiceCritical
	ServiceUnknown
)

func (state ServiceState) String() string {
	switch state {
	case ServiceOK:
		return "ServiceOK"
	case ServiceWarning:
		return "ServiceWarning"
	case ServiceCritical:
		return "ServiceCritical"
	}
	return "ServiceUnknown"
}

// UnmarshalJSON unmarshals service attributes into more meaningful Service field types.
func (s *Service) UnmarshalJSON(data []byte) error {
	type alias Service
	aux := &struct {
		Acknowledgement int
		LastCheck       float64 `json:"last_check"`
		*alias
	}{
		alias: (*alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Acknowledgement != 0 {
		s.Acknowledgement = true
	}
	s.LastCheck = time.Unix(int64(aux.LastCheck), 0)
	return nil
}

func (s Service) Host() string {
	return strings.SplitN(s.Name, "!", 2)[0]
}

func (cr CheckResult) RawCommand() string {
	switch v := cr.Command.(type) {
	case string:
		return v
	case []interface{}:
		var cmd []string
		for i := range v {
			if arg, ok := v[i].(string); ok {
				cmd = append(cmd, arg)
			}
		}
		return strings.Join(cmd, " ")
	}
	return "no command"
}
