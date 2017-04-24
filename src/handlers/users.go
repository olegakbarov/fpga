package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/olegakbarov/io.confs.api/src/db"
)

func GetToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// compare username/pass
	recs, err := db.Read()

	if err != nil {
		log.Fatal("Error quering the db- " + err.Error())
		w.WriteHeader(500)
		return
	}

	// hash, _ := HashPassword(password) // ignore error for the sake of simplicity
	// match := CheckPasswordHash(password, hash)

	res := Envelope{
		Result: "OK",
		Data:   recs,
	}

	data, err := json.Marshal(res)

	if err != nil {
		log.Fatal("Failed marshaling json" + err.Error())
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(data)
}
