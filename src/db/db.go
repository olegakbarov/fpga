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
