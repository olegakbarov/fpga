package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/olegakbarov/io.confs.core/src/db"
	"github.com/olegakbarov/io.confs.core/src/handlers"
)

func main() {
	db.InitDB()

	router := httprouter.New()

	router.GET("/api/v1/conf", handlers.GetAll)
	router.GET("/api/v1/conf/:id", handlers.GetOne)
	router.POST("/api/v1/conf", handlers.Create)
	router.DELETE("/api/v1/conf/:id", handlers.DeleteOne)
	router.PUT("/api/v1/conf/:id", handlers.Edit)

	router.POST("/api/v1/conf/signup", handlers.HandleSignup)
	router.POST("/api/v1/conf/login", handlers.HandleLogin)
	router.POST("/api/v1/conf/logout", auth.CheckAuth(handlers.HandleLogout))

	log.Fatal(http.ListenAndServe(":9999", router))
}
