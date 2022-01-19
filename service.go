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
	Name            string `json:"__name"`
	Groups          []string
	State           ServiceState
	StateType       StateType   `json:"state_type"`
	CheckCommand    string      `json:"check_command"`
	DisplayName     string      `json:"display_name"`
	LastCheckResult CheckResult `json:"last_check_result"`
	Acknowledgement bool
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

func (s ServiceState) String() string {
	switch s {
	case ServiceOK:
		return "ServiceOK"
	case ServiceWarning:
		return "ServiceWarning"
	case ServiceCritical:
		return "ServiceCritical"
	case ServiceUnknown:
		return "ServiceUnknown"
	}
	return "unhandled service state"
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
