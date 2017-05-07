package users

import (
	"time"

	"github.com/guregu/null"
)

type User struct {
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	Locale    null.String `json:"locale"`
	City      null.String `json:"city"`
	Userpic   null.String `json:"userpic"`
	Settings  PropertyMap `json:"settings"`
}

type RawUser struct {
	User
	PasswordHash string    `json:"password_hash"`
	Deleted      bool      `json:"deleted"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
}

func (r RawUser) PublicFields() Conf {
	return r.User
}
