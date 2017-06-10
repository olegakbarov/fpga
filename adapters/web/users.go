package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/olegakbarov/io.confs.core/src/db"
	"github.com/olegakbarov/io.confs.core/src/utils"
)

var SECRET string = os.Getenv("SECRET")

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var rec db.User

	err := decoder.Decode(&rec)
	fmt.Printf("%s\n", &rec)

	user, err := users.GetOne(&rec.id)

	token, err := GenerateToken()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(401)
	}

	res := Envelope{
		Result: "OK",
		Data:   token,
	}

	data, err := json.Marshal(res)

	if err != nil {
		log.Fatal("Failed marshaling json" + err.Error())
		w.WriteHeader(500)
		return
	}

	utils.SendRespose(w, data)
}

func HandleLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// validate shit
	// hash pwd
	// save to db
	// 200ok
}

func HandleLogout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// validate shit
	// hash pwd
	// save to db
	// 200ok
}
