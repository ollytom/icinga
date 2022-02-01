package icinga_test

import (
	"crypto/tls"
	"errors"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"testing"
	"time"

	"olowe.co/icinga"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func randomHostname(suffix string) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b) + suffix
}

func createTestHosts(c *icinga.Client) ([]icinga.Host, error) {
	hostgroup := icinga.HostGroup{Name: "test", DisplayName: "Test Group"}
	if err := c.CreateHostGroup(hostgroup); err != nil && !errors.Is(err, icinga.ErrExist) {
		return nil, err
	}

	var hosts []icinga.Host
	for i := 0; i < 5; i++ {
		h := icinga.Host{
			Name:         randomHostname(".example.org"),
			CheckCommand: "random",
			Groups:       []string{hostgroup.Name},
		}
		hosts = append(hosts, h)
		if err := c.CreateHost(h); err != nil && !errors.Is(err, icinga.ErrExist) {
			return nil, err
		}
	}
	return hosts, nil
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
	for i := range a {
		if a[i] != b[i] {
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

	var want, got []string
	hosts, err := createTestHosts(client)
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range hosts {
		want = append(want, h.Name)
	}
	defer func() {
		for _, h := range hosts {
			if err := client.DeleteHost(h.Name, true); err != nil {
				t.Log(err)
			}
		}
	}()
	hosts, err = client.Hosts("match(\"*example.org\", host.name)")
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

	h := icinga.Host{
		Name:         randomHostname(".checker.example.com"),
		CheckCommand: "hostalive",
	}
	if err := client.CreateHost(h); err != nil {
		t.Fatal(err)
	}

	s := icinga.Service{
		Name:         h.Name + "!http",
		CheckCommand: "http",
	}
	if err := client.CreateService(s); err != nil {
		t.Fatal(err)
	}
	if err := s.Check(client); err != nil {
		t.Fatal(err)
	}
	s, err = client.LookupService(h.Name + "!http")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", s)
}

func TestCheckHostGroup(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Skipf("no local test icinga? got: %v", err)
	}
	hosts, err := createTestHosts(client)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		for _, h := range hosts {
			if err := client.DeleteHost(h.Name, true); err != nil {
				t.Error(err)
			}
		}
	}()
	hostgroup, err := client.LookupHostGroup("test")
	if err != nil {
		t.Fatal(err)
	}
	if err := hostgroup.Check(client); err != nil {
		t.Fatal(err)
	}
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

func TestNonExistentService(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Skipf("no local test icinga? got: %v", err)
	}

	filter := `match("blablabla", service.name)`
	service, err := client.Services(filter)
	if err == nil {
		t.Fail()
	}
	t.Log(err)
	t.Logf("%+v", service)
}
