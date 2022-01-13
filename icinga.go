// package icinga provides a client to the Icinga2 HTTP API.
//
// A Client manages interaction with an Icinga2 server.
// It is created using Dial:
//
//	client, err := icinga.Dial("icinga.example.com:5665", "icinga", "secret", http.DefaultClient)
//	if err != nil {
//		// handle error
//	}
//
// Icinga2 servers in the wild often serve self-signed certificates which fail
// verification by Go's tls client. To ignore the errors, Dial the server with a
// modified http.Client:
//
//	t := http.DefaultTransport.(*http.Transport)
//	t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
//	c := http.DefaultClient
//	c.Transport = t
//	client, err := icinga.Dial(addr, user, pass, c)
//	if err != nil {
//		// handle error
//	}
//
// Methods on Client provide API actions like looking up users and creating
// hosts:
//
//	user, err := client.LookupUser("oliver")
//	if err != nil {
//		// handle error
//	}
//	host := Host{
//		Name: "myserver.example.com",
//		CheckCommand: "hostalive"
//		Address: "192.0.2.1"
//		Address6: "2001:db8::1"
//	}
//	if err := client.CreateHost(host); err != nil {
//		// handle error
//	}
//
// Since Client wraps http.Client, exported methods of http.Client such
// as Get and PostForm can be used to implement any extra functionality
// not provided by this package. For example:
//
//	resp, err := client.PostForm("https://icinga.example.com:5665", data)
//	if err != nil {
//		// handle error
//	}
//
package icinga

import (
	"errors"
	"net/http"
)

// A Client represents a client connection to the Icinga2 HTTP API.
// It should be created using Dial.
// Since Client wraps http.Client, exported methods such as Get and
// PostForm can be used to implement any functionality not provided by
// methods of Client.
type Client struct {
	addr     string
	username string
	password string
	*http.Client
}

var ErrNotExist = errors.New("object does not exist")
var ErrExist = errors.New("object already exists")
var ErrNoMatch = errors.New("no object matches filter")

// Dial returns a new Client connected to the Icinga2 server at addr.
// The recommended value for client is http.DefaultClient.
// But it may also be a modified client which, for example,
// skips TLS certificate verification.
func Dial(addr, username, password string, client *http.Client) (*Client, error) {
	c := &Client{addr, username, password, client}
	if _, err := c.Permissions(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) Permissions() (response, error) {
	resp, err := c.get("", "")
	if err != nil {
		return response{}, err
	}
	if resp.StatusCode == http.StatusOK {
		return response{}, nil
	}
	return response{}, errors.New(resp.Status)
}
