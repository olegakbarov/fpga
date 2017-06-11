package web

import (
	"net/http"

	"github.com/alioygur/gores"
	"github.com/olegakbarov/io.confs.core/core"
)

type (
	errHandlerFunc func(http.ResponseWriter, *http.Request) error

	errResponse struct {
		Code     errCode `json:"code"`
		HTTPCode int     `json:"httpCode"`
		Error    string  `json:"error"`
		Inner    string  `json:"inner,omitempty"`
	}
)

func (h errHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		h.handle(w, err)
	}
}

func (h errHandlerFunc) handle(w http.ResponseWriter, err error) {
	httpErr := func() *webErr {

		if err == core.ErrNoRows {
			return newWebErr(noRowsErrCode, http.StatusNotFound, err)
		}

		switch t := err.(type) {
		case *webErr:
			return t
		case *core.ValidationErr:
			return newWebErr(validationErrCode, http.StatusBadRequest, err)
		}
		// default
		return newWebErr(unknownErrCode, http.StatusInternalServerError, err)
	}()

	errRes := errResponse{Code: httpErr.Code, HTTPCode: httpErr.HTTPCode, Error: httpErr.Error()}
	if httpErr.inner != nil {
		errRes.Inner = httpErr.inner.Error()
	}
	gores.JSON(w, httpErr.HTTPCode, errRes)
}
