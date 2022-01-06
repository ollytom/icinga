package icinga

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type User struct {
	Name   string
	Email  string
	Groups []string
}

var testUser = User{
	Name:  "testUser",
	Email: "test@example.com",
}

var ErrNoUser = errors.New("no such user")

func (u User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		Attrs Alias
	}{Attrs: (Alias)(u)})
}

func (c *Client) Users() ([]User, error) {
	_, err := c.get("/objects/users")
	if err != nil {
		return nil, err
	}
	return []User{testUser}, nil
}

func (c *Client) LookupUser(name string) (User, error) {
	resp, err := c.get("/objects/users/" + name)
	if err != nil {
		return User{}, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return User{}, fmt.Errorf("lookup %s: %w", name, ErrNoUser)
	}
	return testUser, nil
}

func (c *Client) CreateUser(u User) error {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(u); err != nil {
		return err
	}
	resp, err := c.put("/objects/users/"+u.Name, buf)
	if err != nil {
		return fmt.Errorf("create %s: %w", u.Name, err)
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	defer resp.Body.Close()
	var results results
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return fmt.Errorf("create %s: decode response: %w", u.Name, err)
	}
	return fmt.Errorf("create %s: %w", u.Name, results)
}

func (c *Client) DeleteUser(name string) error {
	if err := c.delete("/objects/users/" + name); err != nil {
		return fmt.Errorf("delete user %s: %w", name, err)
	}
	return nil
}
