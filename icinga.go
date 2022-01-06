// package icinga provides a client to the Icinga2 HTTP API.
//
// A Client manages interaction with an Icinga2 server.
// It is created using Dial:
//
//	client, err := icinga.Dial("icinga.example.com:5665", "icinga", "secret", http.DefaultClient)
//	if err != nil {
//		// handle error
//	}
//	host, err := icinga.LookupHost("myserver.example.com")
//	if err != nil {
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
	"fmt"
	"net/http"
)

// A Client represents a client connection to the Icinga2 HTTP API.
// It should be created using Dial.
// Since Client wraps http.Client, standard methods such as Get and
// PostForm can be used to implement any functionality not provided by
// methods of Client.
type Client struct {
	addr     string
	username string
	password string
	*http.Client
}

// Dial returns a new Client connected to the Icinga2 server at addr.
// The recommended value for client is http.DefaultClient.
// But it may also be a modified client which, for example,
// skips TLS certificate verification.
func Dial(addr, username, password string, client *http.Client) (*Client, error) {
	c := &Client{addr, username, password, client}
	if _, err := c.Status(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) Status() (*http.Response, error) {
	resp, err := c.get("/status")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("status %s", resp.Status)
	}
	return resp, nil
}
