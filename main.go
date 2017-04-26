package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/olegakbarov/io.confs.api/src/db"
	"github.com/olegakbarov/io.confs.api/src/handlers"
)

func main() {
	db.InitDB()

	router := httprouter.New()

	router.GET("/", handlers.Index)
	router.GET("/api/v1/conf", handlers.GetAll)
	router.GET("/api/v1/conf/:id", handlers.GetById)
	router.POST("/api/v1/conf", handlers.Add)
	router.DELETE("/api/v1/conf/:id", handlers.DeleteOne)

	router.PUT("/api/v1/conf/:id", handlers.NotImplemented)

	log.Fatal(http.ListenAndServe(":8080", router))
}
