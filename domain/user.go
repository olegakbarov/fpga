package domain

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/guregu/null"
)

type contextKey string

type UserPublicFields struct {
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	Locale    null.String `json:"locale"`
	City      null.String `json:"city"`
	Userpic   null.String `json:"userpic"`
	Settings  PropertyMap `json:"settings"`
}

type User struct {
	Model        `db:",inline"`
	PasswordHash string `json:"password_hash"`
	Confirmed    *bool  `json:"-"`
	UserPublicFields
}

type ReferenceFields struct {
	ID        uint   `json:"id" db:"id,omitempty" gorm:"primary_key"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var userContextKey contextKey = "user"

func (u User) PublicFields() UserPublicFields {
	return u.UserPublicFields
}

func (u User) ReferenceName() ReferenceFields {
	return ReferenceFields{
		u.ID,
		u.FirstName,
		u.LastName,
	}
}

func (u *User) SetPassword(p string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	u.PasswordHash = string(hashedPassword)
}

func (u *User) IsCredentialsVerified(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}

func (u *User) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userContextKey).(*User)
	return u, ok
}

func UserMustFromContext(ctx context.Context) *User {
	u, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		panic("user can't get from request's context")
	}
	return u
}
