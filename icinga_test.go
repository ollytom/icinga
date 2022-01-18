package icinga_test

import (
	"crypto/tls"
	"errors"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"testing"

	"olowe.co/icinga"
)

func randomHostname() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b) + ".example.org"
}

func newTestClient() (*icinga.Client, error) {
	tp := http.DefaultTransport.(*http.Transport)
	tp.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c := http.DefaultClient
	c.Transport = tp
	return icinga.Dial("127.0.0.1:5665", "root", "icinga", c)
}

func compareStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestFilter(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Skipf("no local test icinga? got: %v", err)
	}

	hostgroup := icinga.HostGroup{Name: "examples", DisplayName: "Test Group"}
	if err := client.CreateHostGroup(hostgroup); err != nil {
		t.Error(err)
	}
	hostgroup, err = client.LookupHostGroup(hostgroup.Name)
	if err != nil {
		t.Error(err)
	}
	defer client.DeleteHostGroup(hostgroup.Name, false)

	var want, got []string
	for i := 0; i < 5; i++ {
		h := icinga.Host{
			Name:         randomHostname(),
			CheckCommand: "hostalive",
			Groups:       []string{hostgroup.Name},
		}
		want = append(want, h.Name)
		if err := client.CreateHost(h); err != nil {
			if !errors.Is(err, icinga.ErrExist) {
				t.Error(err)
			}
			continue
		}
		t.Logf("created host %s", h.Name)
	}
	defer func() {
		for _, name := range want {
			if err := client.DeleteHost(name, false); err != nil {
				t.Log(err)
			}
		}
	}()
	hosts, err := client.Hosts("match(\"*example.org\", host.name)")
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range hosts {
		got = append(got, h.Name)
	}
	sort.Strings(want)
	sort.Strings(got)
	if !compareStringSlice(want, got) {
		t.Fail()
	}
	t.Logf("want %+v got %+v", want, got)
}

func TestUserRoundTrip(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Skipf("no local test icinga? got: %v", err)
	}
	want := icinga.User{Name: "olly", Email: "olly@example.com", Groups: []string{}}
	if err := client.CreateUser(want); err != nil && !errors.Is(err, icinga.ErrExist) {
		t.Fatal(err)
	}
	defer func() {
		if err := client.DeleteUser(want.Name, false); err != nil {
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

func TestChecker(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Skipf("no local test icinga? got: %v", err)
	}

	s := icinga.Service{Name: "9p.io!http"}
	if err := s.Check(client); err != nil {
		t.Fatal(err)
	}
	s, err = client.LookupService("9p.io!http")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", s)
}

func TestCreateService(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Skipf("no local test icinga? got: %v", err)
	}

	h := icinga.Host{
		Name:         "example.com",
		Address:      "example.com",
		CheckCommand: "dummy",
		DisplayName:  "RFC 2606 example host",
	}
	if err := client.CreateHost(h); err != nil {
		t.Error(err)
	}
	defer client.DeleteHost(h.Name, true)
	s := icinga.Service{
		Name:         h.Name + "!http",
		CheckCommand: "http",
		DisplayName:  "RFC 2606 example website",
	}
	if err := client.CreateService(s); err != nil {
		t.Error(err)
	}
}
