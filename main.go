package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/olegakbarov/api.confio/db"
	"github.com/olegakbarov/api.confio/handlers"
)

func main() {
	db.InitDB()
	// defer db.Close()

	router := httprouter.New()
	router.GET("/api/v1/conf", handlers.GetAllConfs)
	router.GET("/api/v1/conf/:id", handlers.GetById)
	// router.POST("/api/v1/conf", handlers.AddConf)
	// router.PUT("/api/v1/conf/:id", handlers.UpdateConf)
	// router.DELETE("/api/v1/conf/:id", handlers.DeleteConf)

	log.Fatal(http.ListenAndServe("localhost:1337", router))
}
