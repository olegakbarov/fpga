package db

import (
	"database/sql"
	"fmt"
	"log"
)

func Read() ([]Conf, error) {
	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT * FROM confs ORDER BY id")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rs = make([]RawConf, 0)

	var rec RawConf
	for rows.Next() {
		if err := rows.Scan(
			&rec.Id,
			&rec.Title,
			&rec.Added_by,
			&rec.Start_date,
			&rec.End_date,
			&rec.Description,
			&rec.Picture,
			&rec.Country,
			&rec.City,
			&rec.Address,
			&rec.Category,
			&rec.Tickets_available,
			&rec.Discount_program,
			&rec.Min_price,
			&rec.Max_price,
			&rec.Facebook,
			&rec.Youtube,
			&rec.Twitter,
			&rec.Details,
			&rec.Verified,
			&rec.Deleted,
			&rec.Created_at,
			&rec.Updated_at); err != nil {
			fmt.Printf("%v\n", err)
			return nil, err
		}
		rs = append(rs, rec)
	}

	var pc = make([]Conf, 0)

	for _, item := range rs {
		pc = append(pc, item.PublicFields())
	}

	return pc, nil
}

func ReadOne(id string) (Conf, error) {
	var rec RawConf

	err := db.QueryRow("SELECT * FROM confs WHERE id=$1 ORDER BY id", id).Scan(
		&rec.Id,
		&rec.Title,
		&rec.Added_by,
		&rec.Start_date,
		&rec.End_date,
		&rec.Description,
		&rec.Picture,
		&rec.Country,
		&rec.City,
		&rec.Address,
		&rec.Category,
		&rec.Tickets_available,
		&rec.Discount_program,
		&rec.Min_price,
		&rec.Max_price,
		&rec.Facebook,
		&rec.Youtube,
		&rec.Twitter,
		&rec.Details,
		&rec.Verified,
		&rec.Deleted,
		&rec.Created_at,
		&rec.Updated_at,
	)

	if err != nil {
		log.Fatal(err)
		return Conf{}, err
	} else {
		public := rec.PublicFields()

		return public, nil
	}
}

func EditOne(conf *Conf) (sql.Result, error) {
	stmt, err := db.Prepare("UPDATE confs SET title = $1, start_date = $2, end_date = $3, description = $4, picture = $5, country = $6, city = $7, address = $8, category = $9, tickets_available = $10, discount_program = $11, min_price = $12, max_price = $13, facebook = $14, youtube = $15, twitter = $16, details = $17, id = $18, added_by = $19 WHERE id = $18")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(
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
		&conf.Id,
		&conf.Added_by,
	)

	return res, err
}

func Remove(id string) (sql.Result, error) {
	return db.Exec("DELETE FROM confs WHERE id=$1", id)
}

func Insert(conf *Conf) (sql.Result, error) {
	// TODO! remove hardcoded get_userid()

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
