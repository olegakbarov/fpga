package web

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/olegakbarov/io.confs.core/core"
)

func NewWebAdapter(f core.Factory) http.Handler {
	router := httprouter.New()

	base := alice.New(newSetUserMid(f.NewUser()))
	authRequired := base.Append(newAuthRequiredMid)
	// adminOnly := authRequired.Append(newAdminOnlyMid)

	user := newUser(f)

	router.POST("/api/v1/conf/signup", base.Then(errHandlerFunc(user.signup)))
	router.POST("/api/v1/conf/login", base.Then(errHandlerFunc(user.login)))
	router.POST("/api/v1/conf/logout", base.Then(errHandlerFunc(user.logout)))

	// router.GET("/api/v1/conf", handlers.GetAll)
	// router.GET("/api/v1/conf/:id", handlers.GetOne)
	// router.POST("/api/v1/conf", handlers.Create)
	// router.DELETE("/api/v1/conf/:id", handlers.DeleteOne)
	// router.PUT("/api/v1/conf/:id", handlers.Edit)

	return router
}
