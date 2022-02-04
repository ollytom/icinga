package icinga_test

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"path"
	"reflect"
	"sort"
	"strings"
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

func randomTestAddr() string { return fmt.Sprintf("192.0.2.%d", rand.Intn(254)) }

func randomHosts(n int, suffix string) []icinga.Host {
	var hosts []icinga.Host
	for i := 0; i < n; i++ {
		h := icinga.Host{
			Name:         randomHostname(suffix),
			CheckCommand: "random",
			Groups:       []string{"example"},
			Address:      randomTestAddr(),
		}
		hosts = append(hosts, h)
	}
	return hosts
}

func newTestClient(t *testing.T) *icinga.Client {
	tp := http.DefaultTransport.(*http.Transport)
	tp.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c := &http.Client{Transport: tp}
	if client, err := icinga.Dial("::1:5665", "icinga", "icinga", c); err == nil {
		return client
	}
	client, err := icinga.Dial("127.0.0.1:5665", "icinga", "icinga", c)
	if err == nil {
		return client
	}
	t.Skipf("cannot dial local icinga: %v", err)
	return nil
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
	client := newTestClient(t)
	var want, got []string // host names

	hosts := randomHosts(10, "example.org")
	for i := range hosts {
		if err := client.CreateHost(hosts[i]); err != nil {
			t.Fatal(err)
		}
		want = append(want, hosts[i].Name)
	}
	t.Cleanup(func() {
		for i := range hosts {
			if err := client.DeleteHost(hosts[i].Name, true); err != nil {
				t.Error(err)
			}
		}
	})

	filter := `match("*example.org", host.name)`
	hosts, err := client.Hosts(filter)
	if err != nil {
		t.Fatal(err)
	}
	for i := range hosts {
		got = append(got, hosts[i].Name)
	}

	sort.Strings(want)
	sort.Strings(got)
	if !compareStringSlice(want, got) {
		t.Error("want", want, "got", got)
	}
}

func TestUserRoundTrip(t *testing.T) {
	client := newTestClient(t)
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
	client := newTestClient(t)
	h := randomHosts(1, ".checker.example")[0]
	if err := client.CreateHost(h); err != nil {
		t.Fatal(err)
	}
	defer client.DeleteHost(h.Name, true)
	svc := icinga.Service{
		Name:         h.Name + "!http",
		CheckCommand: "http",
	}
	if err := svc.Check(client); err == nil {
		t.Error("nil error checking non-existent service")
	}
	if err := client.CreateService(svc); err != nil {
		t.Fatal(err)
	}
	if err := svc.Check(client); err != nil {
		t.Error(err)
	}
	if err := client.DeleteService(svc.Name, false); err != nil {
		t.Error(err)
	}
}

func TestCheckHostGroup(t *testing.T) {
	client := newTestClient(t)
	hostgroup := icinga.HostGroup{Name: "test", DisplayName: "Test Group"}
	if err := client.CreateHostGroup(hostgroup); err != nil && !errors.Is(err, icinga.ErrExist) {
		t.Fatal(err)
	}
	defer client.DeleteHostGroup(hostgroup.Name, false)
	hostgroup, err := client.LookupHostGroup(hostgroup.Name)
	if err != nil {
		t.Fatal(err)
	}
	hosts := randomHosts(10, "example.org")
	for _, h := range hosts {
		h.Groups = []string{hostgroup.Name}
		if err := client.CreateHost(h); err != nil {
			t.Fatal(err)
		}
		defer client.DeleteHost(h.Name, false)
	}
	if err := hostgroup.Check(client); err != nil {
		t.Fatal(err)
	}
}

func TestNonExistentService(t *testing.T) {
	client := newTestClient(t)
	filter := `match("blablabla", service.name)`
	service, err := client.Services(filter)
	if err == nil {
		t.Error("non-nil error TODO")
		t.Log(service)
	}
}

type fakeServer struct {
	objects map[string]attributes
}

func newFakeServer() *httptest.Server {
	return httptest.NewTLSServer(&fakeServer{objects: make(map[string]attributes)})
}

// Returns an error message in the same format as returned by the Icinga2 API.
func jsonError(err error) string {
	return fmt.Sprintf("{ %q: %q }", "status", err.Error())
}

var notFoundResponse string = `{
    "error": 404,
    "status": "No objects found."
}`

var alreadyExistsResponse string = `
{
    "results": [
        {
            "code": 500,
            "errors": [
                "Object already exists."
            ],
            "status": "Object could not be created."
        }
    ]
}`

func (srv *fakeServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.RawQuery != "" {
		http.Error(w, jsonError(errors.New("query parameters unimplemented")), http.StatusBadRequest)
		return
	}

	switch {
	case path.Base(req.URL.Path) == "v1":
		srv.Permissions(w)
		return
	case strings.HasPrefix(req.URL.Path, "/v1/objects"):
		srv.ObjectsHandler(w, req)
		return
	}
	http.Error(w, jsonError(errors.New(req.URL.Path+" unimplemented")), http.StatusNotFound)
}

func (f *fakeServer) Permissions(w http.ResponseWriter) {
	fmt.Fprint(w, `{"results": [{
		"info": "Fake Icinga2 server",
		"permissions": ["*"],
		"user": "icinga",
		"version": "fake"
	}]}`)

}

type apiResponse struct {
	Results []apiResult `json:"results"`
	Status  string      `json:"status,omitempty"`
}

type apiResult struct {
	Name  string     `json:"name"`
	Type  string     `json:"type"`
	Attrs attributes `json:"attrs"`
}

// attributes represent configuration object attributes
type attributes map[string]interface{}

// objType returns the icinga2 object type name from an API request path.
// For example from "objects/services/test" the type name is "Service".
func objType(path string) string {
	var t string
	a := strings.Split(path, "/")
	for i := range a {
		if a[i] == "objects" {
			t = a[i+1] // services
		}
	}
	return strings.TrimSuffix(strings.Title(t), "s") // Services to Service
}

func (srv *fakeServer) ObjectsHandler(w http.ResponseWriter, req *http.Request) {
	name := strings.TrimPrefix(req.URL.Path, "/v1/")
	switch req.Method {
	case http.MethodPut:
		if _, ok := srv.objects[name]; ok {
			http.Error(w, alreadyExistsResponse, http.StatusInternalServerError)
			return
		}
		srv.CreateObject(w, req)
	case http.MethodGet:
		srv.GetObject(w, req)
	case http.MethodDelete:
		if _, ok := srv.objects[name]; !ok {
			http.Error(w, notFoundResponse, http.StatusNotFound)
			return
		}
		delete(srv.objects, name)
	default:
		err := fmt.Errorf("%s unimplemented", req.Method)
		http.Error(w, jsonError(err), http.StatusMethodNotAllowed)
	}
}

func (srv *fakeServer) GetObject(w http.ResponseWriter, req *http.Request) {
	name := strings.TrimPrefix(req.URL.Path, "/v1/")
	attrs, ok := srv.objects[name]
	if !ok {
		http.Error(w, notFoundResponse, http.StatusNotFound)
		return
	}
	resp := apiResponse{
		Results: []apiResult{
			apiResult{
				Name:  path.Base(req.URL.Path),
				Type:  objType(req.URL.Path),
				Attrs: attrs,
			},
		},
	}
	json.NewEncoder(w).Encode(&resp)
}

func (srv *fakeServer) CreateObject(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	m := make(map[string]attributes)
	if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
		panic(err)
	}
	name := strings.TrimPrefix(req.URL.Path, "/v1/")
	srv.objects[name] = m["attrs"]
}

func TestDuplicateCreateDelete(t *testing.T) {
	srv := newFakeServer()
	defer srv.Close()
	client, err := icinga.Dial(srv.Listener.Addr().String(), "root", "icinga", srv.Client())
	if err != nil {
		t.Fatal(err)
	}

	host := randomHosts(1, ".example.org")[0]
	if err := client.CreateHost(host); err != nil {
		t.Fatal(err)
	}
	if err := client.CreateHost(host); !errors.Is(err, icinga.ErrExist) {
		t.Errorf("want %s got %v", icinga.ErrExist, err)
	}
	host, err = client.LookupHost(host.Name)
	if err != nil {
		t.Error(err)
	}
	if err := client.DeleteHost(host.Name, false); err != nil {
		t.Error(err)
	}
	if err := client.DeleteHost(host.Name, false); !errors.Is(err, icinga.ErrNotExist) {
		t.Errorf("want icinga.ErrNotExist got %s", err)
	}
	_, err = client.LookupHost(host.Name)
	if !errors.Is(err, icinga.ErrNotExist) {
		t.Errorf("want icinga.ErrNotExist got %s", err)
	}
	if err := client.CreateHost(host); err != nil {
		t.Error(err)
	}
	host, err = client.LookupHost(host.Name)
	if err != nil {
		t.Error(err)
	}
}
