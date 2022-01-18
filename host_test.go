package icinga

import (
	"errors"
	"math/rand"
	"sort"
	"testing"
)

func randomHostname() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b) + ".example.org"
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

	hostgroup := HostGroup{Name: "examples", DisplayName: "Test Group"}
	if err := client.CreateHostGroup(hostgroup); err != nil {
		t.Error(err)
	}
	hostgroup, err = client.LookupHostGroup(hostgroup.Name)
	if err != nil {
		t.Error(err)
	}
	defer client.DeleteHostGroup(hostgroup.Name)

	var want, got []string
	for i := 0; i < 5; i++ {
		h := Host{
			Name:         randomHostname(),
			CheckCommand: "hostalive",
			Groups:       []string{hostgroup.Name},
		}
		want = append(want, h.Name)
		if err := client.CreateHost(h); err != nil {
			if !errors.Is(err, ErrExist) {
				t.Error(err)
			}
			continue
		}
		t.Logf("created host %s", h.Name)
	}
	defer func() {
		for _, name := range want {
			if err := client.DeleteHost(name); err != nil {
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
