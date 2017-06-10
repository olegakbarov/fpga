package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/guregu/null"
)

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
	UserPublicFields
	PasswordHash string    `json:"password_hash"`
	Deleted      bool      `json:"deleted"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
}

func (u User) PublicFields() UserPublicFields {
	return u.UserPublicFields
}

func (u *User) SetPassword(p string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	u.PasswordHash = string(hashedPassword)
}

// compares the givem password with password of user
func (u *User) IsCredentialsVerified(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}

// func (u *User) NewContext(ctx context.Context) context.Context {
//     return context.WithValue(ctx, userContextKey, u)
// }

// // UserFromContext gets user from context
// func UserFromContext(ctx context.Context) (*User, bool) {
//     u, ok := ctx.Value(userContextKey).(*User)
//     return u, ok
// }

// // UserMustFromContext gets user from context. if can't make panic
// func UserMustFromContext(ctx context.Context) *User {
//     u, ok := ctx.Value(userContextKey).(*User)
//     if !ok {
//         panic("user can't get from request's context")
//     }
//     return u
// }
