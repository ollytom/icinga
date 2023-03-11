// Command checkweb is a web application for...
package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path"

	"olowe.co/icinga"
)

type server struct {
	tmpl struct {
		services *template.Template
		service  *template.Template
		hosts    *template.Template
		host     *template.Template
		root     *template.Template
	}
	client *icinga.Client
}

func (srv *server) servicesHandler(w http.ResponseWriter, req *http.Request) {
	filter := req.URL.Query().Get("filter")
	services, err := srv.client.Services(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if err := srv.tmpl.services.Execute(w, services); err != nil {
		log.Println("render services:", err)
	}
}

func (srv *server) hostsHandler(w http.ResponseWriter, req *http.Request) {
	filter := req.URL.Query().Get("filter")
	hosts, err := srv.client.Hosts(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if err := srv.tmpl.hosts.Execute(w, hosts); err != nil {
		log.Println("render hosts:", err)
	}
}

func (srv *server) hostHandler(w http.ResponseWriter, req *http.Request) {
	base := path.Base(req.URL.Path)
	name, err := url.PathUnescape(base)
	if err != nil {
		msg := fmt.Sprintf("escape object name: %v", err)
		log.Println(err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	host, err := srv.client.LookupHost(name)
	if err != nil {
		msg := fmt.Sprintf("lookup %s: %v", name, err)
		code := http.StatusInternalServerError
		if errors.Is(err, icinga.ErrNotExist) {
			code = http.StatusNotFound
		}
		http.Error(w, msg, code)
		return
	}

	if err := srv.tmpl.host.Execute(w, host); err != nil {
		log.Println("render hosts:", err)
	}
}

func (srv *server) serviceHandler(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Path
	if req.URL.RawPath != "" {
		p = req.URL.RawPath
	}
	base := path.Base(p)
	name, err := url.PathUnescape(base)
	if err != nil {
		msg := fmt.Sprintf("escape object name: %v", err)
		log.Println(err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	fmt.Println("looking up", name)
	service, err := srv.client.LookupService(name)
	if err != nil {
		msg := fmt.Sprintf("lookup %s: %v", name, err)
		code := http.StatusInternalServerError
		if errors.Is(err, icinga.ErrNotExist) {
			code = http.StatusNotFound
		}
		http.Error(w, msg, code)
		return
	}

	if err := srv.tmpl.service.Execute(w, service); err != nil {
		log.Println("render service:", err)
	}
}

func (srv *server) searchHandler(w http.ResponseWriter, req *http.Request) {
	var u url.URL
	typ := req.URL.Query().Get("type")
	switch typ {
	case "host":
		u.Path = "/objects/hosts"
	case "service":
		u.Path = "/objects/services"
	}
	v := url.Values{}
	v.Set("filter", req.URL.Query().Get("filter"))
	u.RawQuery = v.Encode()
	http.Redirect(w, req, u.String(), http.StatusFound)
}

func (srv *server) rootHandler(w http.ResponseWriter, req *http.Request) {
	if err := srv.tmpl.root.Execute(w, nil); err != nil {
		log.Println(err)
	}
}

var tFuncMap = template.FuncMap{
	"pathescape": url.PathEscape,
}

func newTemplate(root *template.Template, files ...string) (*template.Template, error) {
	t, err := root.Clone()
	if err != nil {
		return nil, err
	}
	return t.ParseFiles(files...)
}

func main() {
	t := http.DefaultTransport.(*http.Transport)
	t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c := http.DefaultClient
	c.Transport = t
	client, err := icinga.Dial("127.0.0.1:5665", "", "", c)
	if err != nil {
		log.Fatal(err)
	}
	srv := &server{
		client: client,
	}

	srv.tmpl.root = template.Must(template.ParseFiles("base.tmpl")).Funcs(tFuncMap)
	srv.tmpl.services = template.Must(newTemplate(srv.tmpl.root, "services.tmpl"))
	srv.tmpl.service = template.Must(newTemplate(srv.tmpl.root, "service.tmpl"))
	srv.tmpl.hosts = template.Must(newTemplate(srv.tmpl.root, "hosts.tmpl"))
	srv.tmpl.host = template.Must(newTemplate(srv.tmpl.root, "host.tmpl"))

	http.HandleFunc("/objects/services", srv.servicesHandler)
	http.HandleFunc("/objects/services/", srv.serviceHandler)
	http.HandleFunc("/objects/hosts", srv.hostsHandler)
	http.HandleFunc("/objects/hosts/", srv.hostHandler)
	http.HandleFunc("/search", srv.searchHandler)
	http.HandleFunc("/", srv.rootHandler)
	log.Fatal(http.ListenAndServe(":6969", nil))
}
