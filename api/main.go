package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"upper.io/db.v2/postgresql"

	"github.com/olegakbarov/io.confs.core/adapters/web"
	"github.com/olegakbarov/io.confs.core/core"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	dbURL, err := parseDBURL(getConnUrl())
	connURL, err := postgresql.ParseURL(dbURL)
	checkErr(err)

	session, err := postgresql.Open(connURL)
	checkErr(err)

	var sf core.StorageFactory
	sf = postgresqlrepo.NewStorage(session)

	var (
		validator  core.Validator
		mailSender core.MailSender
		jwt        core.JWTSignParser
	)

	validator = providers.NewValidator()
	mailSender = providers.NewFakeMail()
	jwt = providers.NewJWT()
	emitter := providers.NewEmitter()

	f := core.New(sf, mailSender, validator, jwt, emitter)

	log.Printf("Server is up on port: %s", os.Getenv("PORT"))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), web.NewWebAdapter(f)); err != nil {
		log.Fatal(err)
	}
}

func getConnUrl() string {
	s := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s user=%s password=%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		"disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"))

	log.Printf("Connection string looks like this: %s", s)

	return s
}
