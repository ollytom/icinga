package icinga

import "encoding/json"

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
