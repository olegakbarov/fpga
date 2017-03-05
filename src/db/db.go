package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/apex/log"
	"github.com/fogcreek/mini"
)

func getConfig() string {
	ctx := log.WithFields(log.Fields{
		"file": "db.go",
		"func": "Accessing db error",
	})

	cfg, err := mini.LoadConfiguration(".config")

	if err != nil {
		ctx.WithError(err).Error("Error loading configuration")
	}

	info := fmt.Sprintf("host=%s port=%s dbname=%s "+
		"sslmode=%s user=%s password=%s ",
		cfg.String("host", "127.0.0.1"),
		cfg.String("port", "5432"),
		cfg.String("dbname", ""),
		cfg.String("sslmode", "disable"),
		cfg.String("user", ""),
		cfg.String("pass", ""))

	ctx.Info("Config:")
	ctx.Info(info)

	return info
}

var db *sql.DB

func InitDB() {
	var err error
	ctx := log.WithFields(log.Fields{
		"file": "db.go",
		"func": "InitDB()",
	})

	db, err = sql.Open("postgres", getConfig())

	if err != nil {
		ctx.WithError(err).Error("Error connecting to database")
	}

	// defer db.Close()
}

func (p PropertyMap) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

func (p *PropertyMap) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}

func Read() ([]Conf, error) {
	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT * FROM confs ORDER BY id")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rs = make([]Conf, 0)

	var rec Conf
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
			&rec.Facebook_account,
			&rec.Youtube_account,
			&rec.Twitter_account,
			&rec.Details, &rec.Speakers,
			&rec.Sponsors, &rec.Verified,
			&rec.Deleted,
			&rec.Created_at,
			&rec.Updated_at); err != nil {
			fmt.Printf("%v\n", err)
			return nil, err
		}
		rs = append(rs, rec)
	}

	return rs, nil
}

func ReadOne(id string) (Conf, error) {
	var rec Conf
	row := db.QueryRow("SELECT * FROM confs WHERE id=$1 ORDER BY id", id)

	fmt.Printf("%v", row)

	return rec, row.Scan(&rec.Id, &rec.Title, &rec.Added_by, &rec.Start_date, &rec.End_date, &rec.Description, &rec.Picture, &rec.Country, &rec.City, &rec.Address, &rec.Category, &rec.Min_price, &rec.Max_price,
		//&rec.Facebook_account,
		//&rec.Youtube_account,
		//&rec.Twitter_account,
		&rec.Tickets_available, &rec.Discount_program, &rec.Details, &rec.Speakers, &rec.Sponsors, &rec.Verified, &rec.Deleted, &rec.Created_at, &rec.Updated_at)
}

func Remove(id string) (sql.Result, error) {
	return db.Exec("DELETE FROM confs WHERE id=$1", id)
}

func Insert(item Conf) (sql.Result, error) {
	// todo insert one by one??
	return db.Exec("INSERT INT confs VALUES (default, $1)",
		item)
}

// func Update(id int, title string) (sql.Result, error) {
// 	return db.Exec("UPDATE phonebook SET title = $1, WHERE id=$2",
// 		title, id)
// }
