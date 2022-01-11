package icinga

import (
	"encoding/json"
	"fmt"
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

func (u User) MarshalJSON() ([]byte, error) {
	type attrs struct {
		Email  string   `json:"email"`
		Groups []string `json:"groups,omitempty"`
	}
	return json.Marshal(&struct {
		Attrs attrs `json:"attrs"`
	}{
		Attrs: attrs{
			Email:  u.Email,
			Groups: u.Groups,
		},
	})
}

func (u User) name() string {
	return u.Name
}

func (u User) path() string {
	return "/objects/users/" + u.Name
}

func (u User) attrs() map[string]interface{} {
	m := make(map[string]interface{})
	m["groups"] = u.Groups
	m["email"] = u.Email
	return m
}

func (c *Client) Users() ([]User, error) {
	objects, err := c.allObjects("/objects/users")
	if err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	var users []User
	for _, o := range objects {
		v, ok := o.(User)
		if !ok {
			return nil, fmt.Errorf("get all users: %T in response", v)
		}
		users = append(users, v)
	}
	return users, nil
}

func (c *Client) LookupUser(name string) (User, error) {
	obj, err := c.lookupObject("/objects/users/" + name)
	if err != nil {
		return User{}, fmt.Errorf("lookup %s: %w", name, err)
	}
	v, ok := obj.(User)
	if !ok {
		return User{}, fmt.Errorf("lookup %s: result type %T is not user", name, v)
	}
	return v, nil
}

// CreateUser creates user.
// An error is returned if the User already exists or on any other error.
func (c *Client) CreateUser(user User) error {
	if err := c.createObject(user); err != nil {
		return fmt.Errorf("create user %s: %w", user.Name, err)
	}
	return nil
}

// DeleteUser deletes the User identified by name.
// ErrNotExist is returned if the User doesn't exist.
func (c *Client) DeleteUser(name string) error {
	if err := c.deleteObject("/objects/users/" + name); err != nil {
		return fmt.Errorf("delete user %s: %w", name, err)
	}
	return nil
}
