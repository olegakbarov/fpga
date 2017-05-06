package users

import (
	"time"

	"github.com/guregu/null"
)

type RawUser struct {
	User
	Verified   bool      `json:"verified"`
	Deleted    bool      `json:"deleted"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type User struct {
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Locale    null.String `json:"locale"`
	City      null.String `json:"city"`
	Userpic   null.String `json:"userpic"`
	Email     string      `json:"email"`
	Settings  PropertyMap `json:"settings"`
}

func (r RawConf) PublicFields() Conf {
	return r.Conf
}
