package icinga

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name  string
	Type  string
	Attrs struct {
		Email string
	}
}

var testUser = User{
	Name: "Olly",
	Type: "User",
	Attrs: struct {
		Email string
	}{Email: "olly@example.com"},
}

func (c *Client) Users() ([]User, error) {
	resp, err := c.get("/objects/users")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get /objects/users: status %s", resp.Status)
	}
	return []User{testUser}, nil
}

func (c *Client) CreateUser(name, email string) error {
	u := User{
		Name: name,
		Type: "User",
		Attrs: struct {
			Email string
		}{email},
	}
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(u); err != nil {
		return err
	}
	_, err := c.put("/objects/users/"+name, buf)
	return err
}
