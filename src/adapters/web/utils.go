package web

import (
	"net/http"
	"os"
)

func SendRespose(w http.ResponseWriter, data []byte) {
	if os.Getenv("SECRET") == "development" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}
