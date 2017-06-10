package users

import (
	"database/sql"
	"log"
)

func Create(user *User) (sql.Result, error) {
	stmt, err := db.Prepare("INSERT INTO users (first_name, last_name, locale, city, userpic, email, settings, verified, deleted, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING id")
	if err != nil {
		log.Fatal(err)
	}

	return stmt.Exec(
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Locale,
		&user.City,
		&user.Userpic,
		&user.Settings,
		&user.PasswordHash,
	)
}

func GetOne(id string) (db.RawUser error) {

}
