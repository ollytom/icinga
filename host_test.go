package icinga

import (
	"crypto/tls"
	"errors"
	"math/rand"
	"net/http"
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
	tp := http.DefaultTransport.(*http.Transport)
	tp.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c := http.DefaultClient
	c.Transport = tp
	client, err := Dial("127.0.0.1:5665", "root", "8eec5ede1673b757", c)
	if err != nil {
		t.Skipf("no local test icinga? got: %v", err)
	}

	var want, got []string
	for i := 0; i < 5; i++ {
		h := Host{Name: randomHostname(), CheckCommand: "hostalive"}
		want = append(want, h.Name)
		if err := client.CreateHost(h); err != nil {
			if !errors.Is(err, ErrExist) {
				t.Error(err)
			}
			continue
		}
		t.Logf("created host %s", h.Name)
	}
	hosts, err := client.FilterHosts("match(\"*example.org\", host.name)")
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
