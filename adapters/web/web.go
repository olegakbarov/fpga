package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/olegakbarov/io.confs.core/core"
)

func NewWebAdapter(f core.Factory) http.Handler {
	r := mux.NewRouter()

	base := alice.New(newSetUserMid(f.NewUser()))
	// authRequired := base.Append(newAuthRequiredMid)
	// adminOnly := authRequired.Append(newAdminOnlyMid)

	user := newUser(f)

	r.Handle("/api/v1/conf/signup", base.Then(errHandlerFunc(user.signup))).Methods("POST")
	r.Handle("/api/v1/conf/login", base.Then(errHandlerFunc(user.login))).Methods("POST")
	// r.Handle("/api/v1/conf/logout", base.Then(errHandlerFunc(user.logout))).Methods("POST")

	// router.GET("/api/v1/conf", handlers.GetAll)
	// router.GET("/api/v1/conf/:id", handlers.GetOne)
	// router.POST("/api/v1/conf", handlers.Create)
	// router.DELETE("/api/v1/conf/:id", handlers.DeleteOne)
	// router.PUT("/api/v1/conf/:id", handlers.Edit)

	return r
}
