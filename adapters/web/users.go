package web

import (
	"net/http"

	"github.com/alioygur/gores"
	"github.com/olegakbarov/io.confs.core/core"
)

type (
	user struct {
		core.User
	}
)

func newUser(f core.Factory) *user {
	return &user{f.NewUser()}
}

func (u *user) signup(w http.ResponseWriter, r *http.Request) error {
	req := new(core.RegisterRequest)
	if err := decodeReq(r, req); err != nil {
		return err
	}

	usr, err := u.Register(req)
	if err != nil {
		if err == core.ErrEmailExists {
			return newWebErr(emailExistsErrCode, http.StatusConflict, err)
		}
		return err
	}

	jwt, err := u.GenToken(usr, core.AuthToken)
	if err != nil {
		return err
	}

	return gores.JSON(w, http.StatusCreated, response{jwt})
}

func (u *user) login(w http.ResponseWriter, r *http.Request) error {
	req := new(core.LoginRequest)
	if err := decodeReq(r, req); err != nil {
		return err
	}

	usr, err := u.Login(req)
	if err != nil {
		switch err {
		case core.ErrWrongCredentials:
			return newWebErr(wrongCredErrCode, http.StatusUnauthorized, err)
		case core.ErrInActiveUser:
			return newWebErr(inactiveUserErrCode, http.StatusUnauthorized, err)
		}
		return err
	}

	jwt, err := u.GenToken(usr, core.AuthToken)
	if err != nil {
		return err
	}

	return gores.JSON(w, http.StatusOK, response{jwt})
}
