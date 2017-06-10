package web

import (
	"github.com/julienschmidt/httprouter"
	"github.com/olegakbarov/io.confs.core/src/handlers"
)

func NewWebAdapter() {
	router := httprouter.New()

	router.POST("/api/v1/conf/signup", handlers.HandleSignup)
	router.POST("/api/v1/conf/login", handlers.HandleLogin)
	router.POST("/api/v1/conf/logout", auth.CheckAuth(handlers.HandleLogout))

	// router.GET("/api/v1/conf", handlers.GetAll)
	// router.GET("/api/v1/conf/:id", handlers.GetOne)
	// router.POST("/api/v1/conf", handlers.Create)
	// router.DELETE("/api/v1/conf/:id", handlers.DeleteOne)
	// router.PUT("/api/v1/conf/:id", handlers.Edit)

	return router
}
