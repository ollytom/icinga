package icinga

import (
	"os"
	"testing"
)

func TestHostUnmarshal(t *testing.T) {
	f, err := os.Open("testdata/objects/hosts/VuS9jZ8u.example.org")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	resp, err := parseResponse(f)
	if err != nil {
		t.Fatal(err)
	}
	host := resp.Results[0].(Host)
	if host.LastCheck.IsZero() {
		t.Error("zero time")
	}
	if !host.Acknowledgement {
		t.Error("should be acknowledged")
	}
	if t.Failed() {
		t.Log(host)
	}
}
