package icinga

import (
	"os"
	"reflect"
	"testing"
)

func TestUser(t *testing.T) {
	want := User{Name: "test", Email: "test@example.com", Groups: []string{}}
	f, err := os.Open("testdata/objects/users/test")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	resp, err := parseResponse(f)
	if err != nil {
		t.Fatal(err)
	}
	got := resp.Results[0].(User)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %+v, got %+v", want, got)
	}
}
