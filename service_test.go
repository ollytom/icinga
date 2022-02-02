package icinga

import (
	"os"
	"reflect"
	"testing"
)

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
	want := Service{
		Name:         "9p.io!http",
		Groups:       []string{},
		State:        ServiceOK,
		StateType:    StateHard,
		CheckCommand: "http",
		DisplayName:  "http",
		LastCheckResult: &CheckResult{
			Output: "HTTP OK: HTTP/1.1 200 OK - 1714 bytes in 1.083 second response time ",
		},
	}
	var got Service
	for _, r := range resp.Results {
		if r.name() == "9p.io!http" {
			got = r.(Service)
		}
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %+v, got %+v", want, got)
	}
}
