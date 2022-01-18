package icinga

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestUser(t *testing.T) {
	want := User{Name: "test", Email: "test@example.com", Groups: []string{}}
	f, err := os.Open("testdata/users.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	resp, err := parseResponse(f)
	if err != nil {
		t.Fatal(err)
	}
	obj, err := objectFromLookup(resp)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := obj.(User)
	if !ok {
		t.Fatalf("want %T, got %T", want, got)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %+v, got %+v", want, got)
	}
}

func TestUserMarshal(t *testing.T) {
	user := &User{Name: "test", Email: "test@example.com", Groups: []string{}}
	want := `{"attrs":{"email":"test@example.com"}}`
	got, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fail()
	}
	t.Logf("want %s, got %s", want, got)
}
