package icinga

import (
	"os"
	"testing"
	"time"
)

// Tests the trickier parts of the custom Unmarshaller functionality.
func TestServiceUnmarshal(t *testing.T) {
	f, err := os.Open("testdata/objects/services/9p.io!http")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	resp, err := parseResponse(f)
	if err != nil {
		t.Fatal(err)
	}
	svc := resp.Results[0].(Service)
	if svc.LastCheck.IsZero() {
		t.Error("zero time")
	}
	if !svc.Acknowledgement {
		t.Error("should be acknowledged")
	}
	if t.Failed() {
		t.Log(svc)
	}
}

func TestServiceMarshalForCreate(t *testing.T) {
	want := `{"attrs":{"check_command":"dummy","display_name":"test"}}`
	service := Service{
		CheckCommand: "dummy",
		DisplayName:  "test",
		LastCheck:    time.Now(),
		LastCheckResult: CheckResult{
			Output:      "xxx",
			CheckSource: "xxx",
			Command:     nil,
		},
	}
	got, err := jsonForCreate(service)
	if err != nil {
		t.Fatal(err)
	}
	if want != string(got) {
		t.Error("not matching", string(got))
	}
}
