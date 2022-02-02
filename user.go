package icinga

// User represents a User object.
// Note that this is different from an ApiUser.
type User struct {
	Name   string   `json:"-"`
	Email  string   `json:"email,omitempty"`
	Groups []string `json:"groups,omitempty"`
}

func (u User) name() string {
	return u.Name
}

func (u User) path() string {
	return "/objects/users/" + u.Name
}
