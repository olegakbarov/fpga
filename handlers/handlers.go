package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/julienschmidt/httprouter"
	"github.com/olegakbarov/api.confio/db"
)

func GetAllConfs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := log.WithFields(log.Fields{
		"file": "handlers.go",
		"func": "GetAllConfs",
	})

	recs, err := db.Read()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	data, err := json.Marshal(recs)
	fmt.Printf("%s\n", data)

	if err != nil {
		err := errors.New("err")
		ctx.WithError(err).Error("Failed marshaling json")
		w.WriteHeader(500)
		return
	}

	ctx.Info("Sending response..")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
	ctx.Info("Done.")
}

func GetById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := log.WithFields(log.Fields{
		"file": "handlers.go",
		"func": "GetById",
	})

	id := ps.ByName("id")
	fmt.Printf("%s\n", id)

	rec, err := db.ReadOne(id)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}

		err := errors.New("err")
		ctx.WithError(err).Error("Sql error")
		w.WriteHeader(500)
		return
	}

	data, err := json.Marshal(rec)
	if err != nil {
		err := errors.New("err")
		ctx.WithError(err).Error("Failed marshaling json")
		w.WriteHeader(500)
		return
	}

	ctx.Info("Sending response..")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
	ctx.Info("Done.")
}

// func AddRecord(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	var rec db.Record
// 	err := json.NewDecoder(r.Body).Decode(&rec)
// 	if err != nil || rec.Title == "" {
// 		w.WriteHeader(400)
// 		return
// 	}
// 	if _, err := db.Insert(rec.Title); err != nil {
// 		w.WriteHeader(500)
// 		return
// 	}
// 	w.WriteHeader(201)
// }
//
// func UpdateRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	id, ok := getID(w, ps)
// 	if !ok {
// 		return
// 	}
// 	var rec Record
// 	err := json.NewDecoder(r.Body).Decode(&rec)
// 	if err != nil || rec.Title == "" {
// 		w.WriteHeader(400)
// 		return
// 	}
// 	res, err := db.Update(id, rec.Title)
// 	if err != nil {
// 		w.WriteHeader(500)
// 		return
// 	}
// 	n, _ := res.RowsAffected()
// 	if n == 0 {
// 		w.WriteHeader(404)
// 		return
// 	}
// 	w.WriteHeader(204)
// }
//
// func DeleteRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	id, ok := getID(w, ps)
// 	if !ok {
// 		return
// 	}
// 	if _, err := db.Remove(id); err != nil {
// 		w.WriteHeader(500)
// 	}
// 	w.WriteHeader(204)
// }
