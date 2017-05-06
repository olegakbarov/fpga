package db

import (
	"database/sql"
	"log"
)

func Insert(user *User) (sql.Result, error) {
	stmt, err := db.Prepare("INSERT INTO confs (title, added_by, start_date, end_date, description, picture, country, city, address, category, tickets_available, discount_program, min_price, max_price, facebook, youtube, twitter, details) values ($1, get_userid('Oleg'), $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING id")
	if err != nil {
		log.Fatal(err)
	}

	return stmt.Exec(
		&conf.Title,
		&conf.Start_date,
		&conf.End_date,
		&conf.Description,
		&conf.Picture,
		&conf.Country,
		&conf.City,
		&conf.Address,
		&conf.Category,
		&conf.Tickets_available,
		&conf.Discount_program,
		&conf.Min_price,
		&conf.Max_price,
		&conf.Facebook,
		&conf.Youtube,
		&conf.Twitter,
		&conf.Details,
	)
}
