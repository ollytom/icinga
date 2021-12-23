package icinga

import (
	"fmt"
	"net/http"
)

type Client struct {
	host       string
	username   string
	password   string
	httpClient *http.Client
}

func Dial(host, username, password string, client *http.Client) (*Client, error) {
	c := &Client{host, username, password, client}
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
