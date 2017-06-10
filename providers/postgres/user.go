package postgres

import "log"

func GetOne(id string) (User error) {
	var rec users.RawUser

	err := db.QueryRow("SELECT * FROM users WHERE id=$1 ORDER BY id", id).Scan(
		&rec.FirstName,
		&rec.LastName,
		&rec.Email,
		&rec.Locale,
		&rec.City,
		&rec.Userpic,
		&rec.Settings,
		&rec.PasswordHash,
		&rec.Deleted,
		&rec.Created_at,
		&rec.Updated_at,
	)

	if err != nil {
		log.Fatal("Error quering the db- " + err.Error())
		w.WriteHeader(500)
		return
	}

	return rec.PublicFields(), nil
}
