package icinga

import (
	"os"
	"reflect"
	"testing"
)

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
