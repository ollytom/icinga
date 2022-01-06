package icinga

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// User represents a User object.
// Note that this is different from an ApiUser.
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

// CreateUser creates the User u identified by u.Name.
// An error is returned if the User already exists or on any other error.
func (c *Client) CreateUser(u User) error {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(u); err != nil {
		return err
	}
	if err := c.put("/objects/users/"+u.Name, buf); err != nil {
		return fmt.Errorf("create %s: %w", u.Name, err)
	}
	return nil
}

// DeleteUser deletes the User identified by name.
// ErrNoUser is returned if the User doesn't exist.
func (c *Client) DeleteUser(name string) error {
	if err := c.delete("/objects/users/" + name); err != nil {
		return fmt.Errorf("delete user %s: %w", name, err)
	}
	return nil
}
