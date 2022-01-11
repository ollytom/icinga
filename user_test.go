package icinga

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
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

func TestUserRoundTrip(t *testing.T) {
	tp := http.DefaultTransport.(*http.Transport)
	tp.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c := http.DefaultClient
	c.Transport = tp
	client, err := Dial("127.0.0.1:5665", "root", "8eec5ede1673b757", c)
	if err != nil {
		t.Skipf("no local test icinga? got: %v", err)
	}

	want := User{Name: "olly", Email: "olly@example.com", Groups: []string{}}
	if err := client.CreateUser(want); err != nil && !errors.Is(err, ErrExist) {
		t.Fatal(err)
	}
	defer func() {
		if err := client.DeleteUser(want.Name); err != nil {
			t.Error(err)
		}
	}()
	got, err := client.LookupUser(want.Name)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %+v, got %+v", want, got)
	}
}
