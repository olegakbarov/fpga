package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"upper.io/db.v2/mysql"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/alioygur/gocart/adapters/web"
	"github.com/alioygur/gocart/engine"
	"github.com/alioygur/gocart/providers"
	mysqlrepo "github.com/alioygur/gocart/providers/mysql"
)

func main() {
	dbURL, err := parseDBURL(os.Getenv("MYSQL_URL"))
	if err != nil {
		log.Fatal(err)
	}
	connURL, err := mysql.ParseURL(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	session, err := mysql.Open(connURL)
	if err != nil {
		log.Fatal(err)
	}

	// Setup storage factory
	var sf engine.StorageFactory
	sf = mysqlrepo.NewStorage(session)

	// Setup service dependencies
	var (
		validator  engine.Validator
		mailSender engine.MailSender
		jwt        engine.JWTSignParser
	)

	validator = providers.NewValidator()
	mailSender = providers.NewFakeMail()
	jwt = providers.NewJWT()
	emitter := providers.NewEmitter()

	f := engine.New(sf, mailSender, validator, jwt, emitter)

	log.Printf("server starting port: %s", os.Getenv("PORT"))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), web.NewWebAdapter(f)); err != nil {
		log.Fatal(err)
	}
}

func parseDBURL(s string) (string, error) {
	durl, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("cannot parse database url, err:%s", err)
	}
	user := durl.User.Username()
	password, _ := durl.User.Password()
	host := durl.Host
	dbname := durl.Path // like: /path

	return fmt.Sprintf("%s:%s@tcp(%s)%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbname), nil
}
