// this package is just for development, and, not intended to be used in production.
package inmem

import (
	"fmt"
	"time"

	"github.com/betelgeuse-7/qa/service/hashpwd"
)

var users []*User

type User struct {
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Username     string     `uri:"username" json:"username"`
	Email        string     `json:"email"`
	Password     string     `json:"-"`
	RegisteredAt *time.Time `json:"registered_at"`
}

func (u *User) Save() error {
	if u := GetUser(u.Username); u != nil {
		return fmt.Errorf("user with username '%s', already exists", u.Username)
	}
	hasher := hashpwd.New(u.Password)
	hasher.HashPwd()
	if err := hasher.Error(); err != nil {
		return err
	}
	u.Password = hasher.Hashed()
	users = append(users, u)
	return nil
}

func GetUser(username string) *User {
	for _, u := range users {
		if u.Username == username {
			return u
		}
	}
	return nil
}
