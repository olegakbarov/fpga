package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func getConfig() string {
	info := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s user=%s password=%s ",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		"disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"))

	log.Printf("Config looks like this: %s", info)

	return info
}

func getID(w http.ResponseWriter, ps httprouter.Params) (int, bool) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		return 0, false
	}
	return id, true
}

var db *sql.DB

func InitDB() {
	var err error

	db, err = sql.Open("postgres", getConfig())

	if err != nil {
		log.Fatal("Error: The data source arguments are not valid - " + err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
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

func EditOne(conf Conf) (sql.Result, error) {
	res, err := db.Exec("UPDATE confs SET id= $1, Added_by = $2 title = $3, start_date = $4, end_data = $5, description = $6, picture = $7, country = $8, city = $9, address = $10, category = $11, tickets_available = $12, discount_program = $13, min_price = $14, max_price = $15, facebook = $16, youtube = $17, twitter = $18, details = $19 WHERE id=$1",
		conf.Id,
		conf.Added_by,
		conf.Title,
		conf.Start_date,
		conf.End_date,
		conf.Description,
		conf.Picture,
		conf.Country,
		conf.City,
		conf.Address,
		conf.Category,
		conf.Tickets_available,
		conf.Discount_program,
		conf.Min_price,
		conf.Max_price,
		conf.Facebook,
		conf.Youtube,
		conf.Twitter,
		conf.Details,
	)

	return res, err
}

func Remove(id string) (sql.Result, error) {
	return db.Exec("DELETE FROM confs WHERE id=$1", id)
}

func Insert(item Conf) (sql.Result, error) {
	// todo insert one by one??
	return db.Exec("INSERT INTO confs VALUES (default, $1)", item)
}
