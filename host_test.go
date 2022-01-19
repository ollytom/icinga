package icinga

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestHostMarshal(t *testing.T) {
	s := `{"attrs":{"address":"192.0.2.1","address6":"2001:db8::","groups":["test"],"check_command":"dummy","display_name":"Example host"}}`
	want := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &want); err != nil {
		t.Fatal(err)
	}

	p, err := json.Marshal(Host{
		Name:         "example.com",
		Address:      "192.0.2.1",
		Address6:     "2001:db8::",
		Groups:       []string{"test"},
		State:        HostDown,
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
		t.Fail()
	}
	t.Log("want", want, "got", got)
}

func TestHostUnmarshal(t *testing.T) {
	want := Host{
		Name:            "example.com",
		Address:         "",
		Groups:          []string{"example"},
		State:           HostDown,
		StateType:       StateSoft,
		CheckCommand:    "hostalive",
		DisplayName:     "example.com",
		Acknowledgement: false,
	}
	f, err := os.Open("testdata/hosts.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	resp, err := parseResponse(f)
	if err != nil {
		t.Fatal(err)
	}
	var got Host
	for _, r := range resp.Results {
		h := r.(Host)
		if h.Name == "example.com" {
			got = h
			break
		}
	}
	if !reflect.DeepEqual(want, got) {
		t.Fail()
	}
	t.Logf("want %+v, got %+v", want, got)
}
