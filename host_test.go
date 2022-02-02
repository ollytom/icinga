package icinga

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestHostMarshal(t *testing.T) {
	b := []byte(`{"attrs":{"address":"192.0.2.1","address6":"2001:db8::","check_command":"dummy","display_name":"Example host","groups":["test"]}}`)
	want := make(map[string]interface{})
	if err := json.Unmarshal(b, &want); err != nil {
		t.Fatal(err)
	}

	p, err := json.Marshal(Host{
		Name:         "example.com",
		Address:      "192.0.2.1",
		Address6:     "2001:db8::",
		Groups:       []string{"test"},
		StateType:    StateSoft,
		CheckCommand: "dummy",
		DisplayName:  "Example host",
	})
	if err != nil {
		t.Fatal(err)
	}
	got := make(map[string]interface{})
	if err := json.Unmarshal(p, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Error("want", want, "got", got)
	}
}

func TestHostUnmarshal(t *testing.T) {
	want := Host{
		Name:            "VuS9jZ8u.example.org",
		Address:         "",
		Groups:          []string{},
		State:           HostDown,
		StateType:       StateSoft,
		CheckCommand:    "hostalive",
		DisplayName:     "VuS9jZ8u.example.org",
		Acknowledgement: false,
	}
	f, err := os.Open("testdata/objects/hosts/VuS9jZ8u.example.org")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	resp, err := parseResponse(f)
	if err != nil {
		t.Fatal(err)
	}
	got := resp.Results[0].(Host)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %+v, got %+v", want, got)
	}
}
