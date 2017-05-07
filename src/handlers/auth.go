package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
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

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		tokens, ok := r.Header["Authorization"]
		if ok && len(tokens) >= 1 {
			token = tokens[0]
			token = strings.TrimPrefix(token, "Bearer ")
		}

		if token == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				msg := fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				return nil, msg
			}
			return a.encryptionKey, nil
		})

		if err != nil {
			http.Error(w, "Error parsing token", http.StatusUnauthorized)
			return
		}

		if parsedToken != nil && parsedToken.Valid {
			context.Set(r, "user", parsedToken)
			next.ServeHTTP(w, r)
			fmt.Println("yay! token is valid")
		}

		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	})
}
