package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/olegakbarov/api.confsio/src/db"
	"github.com/olegakbarov/api.confsio/src/handlers"
)

func main() {
	db.InitDB()

	router := httprouter.New()

	router.GET("/", handlers.Index)
	router.GET("/api/v1/conf", handlers.GetAllConfs)
	router.GET("/api/v1/conf/:id", handlers.GetById)
	router.POST("/api/v1/conf", handlers.AddConf)
	// router.PUT("/api/v1/conf/:id", handlers.UpdateConf)
	router.DELETE("/api/v1/conf/:id", handlers.DeleteConfById)

	log.Fatal(http.ListenAndServe(":8080", router))
}
