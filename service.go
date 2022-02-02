package icinga

import "encoding/json"

func (s Service) name() string {
	return s.Name
}

func (s Service) path() string {
	return "/objects/services/" + s.Name
}

// Service represents a Service object.
type Service struct {
	Name            string   `json:"-"`
	Groups          []string `json:"groups,omitempty"`
	State           ServiceState
	StateType       StateType   `json:"state_type"`
	CheckCommand    string      `json:"check_command"`
	DisplayName     string      `json:"display_name,omitempty"`
	LastCheckResult CheckResult `json:"last_check_result,omitempty"`
	Acknowledgement bool        `json:",omitempty"`
}

type CheckResult struct {
	Output string
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

func (s Service) MarshalJSON() ([]byte, error) {
	attrs := make(map[string]interface{})
	if len(s.Groups) > 0 {
		attrs["groups"] = s.Groups
	}
	attrs["check_command"] = s.CheckCommand
	attrs["display_name"] = s.DisplayName
	jservice := &struct {
		Attrs map[string]interface{} `json:"attrs"`
	}{Attrs: attrs}
	return json.Marshal(jservice)
}

// UnmarshalJSON unmarshals service attributes into more meaningful Service field types.
func (s *Service) UnmarshalJSON(data []byte) error {
	type alias Service
	aux := &struct {
		Acknowledgement int
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
	return nil
}
