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
	router.GET("/api/v1/conf/:id", handlers.GetOne)
	router.POST("/api/v1/conf", handlers.Create)
	router.DELETE("/api/v1/conf/:id", handlers.DeleteOne)
	router.PUT("/api/v1/conf/:id", handlers.Edit)

	router.POST("/api/v1/auth/get_token", handlers.NotImplemented)
	router.POST("/api/v1/conf/login", handlers.NotImplemented)
	router.POST("/api/v1/conf/logout", handlers.NotImplemented)
	router.POST("/api/v1/conf/signup", handlers.NotImplemented)

	log.Fatal(http.ListenAndServe(":9999", router))
}
