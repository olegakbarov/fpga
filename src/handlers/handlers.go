package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/julienschmidt/httprouter"
	"github.com/olegakbarov/api.confsio/src/db"
)

type Envelope struct {
	Result string      `json:"result"`
	Data   interface{} `json:"data"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func GetAllConfs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := log.WithFields(log.Fields{
		"file": "handlers.go",
		"func": "GetAllConfs",
	})

	recs, err := db.Read()
	if err != nil {
		err := errors.New("database error")
		ctx.WithError(err).Error("Check db handler")
		w.WriteHeader(500)
		return
	}

	res := Envelope{
		Result: "Success",
		Data:   recs,
	}

	data, err := json.Marshal(res)

	if err != nil {
		err := errors.New("err")
		ctx.WithError(err).Error("Failed marshaling json")
		w.WriteHeader(500)
		return
	}

	ctx.Info("Sending response..")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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

func DeleteConfById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := log.WithFields(log.Fields{
		"file": "handlers.go",
		"func": "DeleteConfById",
	})

	id := ps.ByName("id")
	fmt.Printf("%s\n", id)

	if _, err := db.Remove(id); err != nil {
		err := errors.New("err")
		ctx.WithError(err).Error("Error deleting msg")
		w.WriteHeader(500)
	}

	w.WriteHeader(204)
}

func AddConf(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := log.WithFields(log.Fields{
		"file": "handlers.go",
		"func": "AddConf",
	})

	decoder := json.NewDecoder(r.Body)
	var rec db.Conf

	err := decoder.Decode(&rec)
	fmt.Printf("%s\n", &rec)

	if err != nil {
		ctx.WithError(err).Error("Error marshaling json")
		w.WriteHeader(400)
		return
	}
	if _, err := db.Insert(rec); err != nil {
		ctx.WithError(err).Error("Error insering into db")
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)
}

//func updaterecord(w http.responsewriter, r *http.request, ps httprouter.params) {
//id, ok := getid(w, ps)
//if !ok {
//return
//}
//var rec record
//err := json.newdecoder(r.body).decode(&rec)
//if err != nil || rec.title == "" {
//w.writeheader(400)
//return
//}
//res, err := db.update(id, rec.title)
//if err != nil {
//w.writeheader(500)
//return
//}
//n, _ := res.rowsaffected()
//if n == 0 {
//w.writeheader(404)
//return
//}
//w.writeheader(204)
//}
