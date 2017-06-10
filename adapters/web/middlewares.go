package web

import (
	"errors"
	"net/http"
	"strings"

	"github.com/olegakbarov/io.confs.core/core"
	"github.com/olegakbarov/io.confs.core/domain"
)

type (
	userFinder interface {
		GetFromAuthToken(tokenStr string) (*domain.User, error)
	}
)

// TODO CORSWrapper
// if os.Getenv("SECRET") == "development" {
//     w.Header().Set("Access-Control-Allow-Origin", "*")
// }

func getToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("invalid authorization token format")
	}

	return authHeaderParts[1], nil
}

func newSetUserMid(uf userFinder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return errHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			tokenStr, err := getToken(r)
			if err != nil {
				return newWebErr(badRequestErrCode, http.StatusBadRequest, err)
			}

			if tokenStr == "" {
				next.ServeHTTP(w, r)
				return nil
			}

			u, err := uf.GetFromAuthToken(tokenStr)
			if err != nil {
				// handle token err
				tokenErr, ok := err.(*core.TokenErr)
				if ok {
					if tokenErr.Expired() {
						return newWebErr(expiredTokenErrCode, http.StatusUnauthorized, err)
					}
					return newWebErr(invalidTokenErrCode, http.StatusUnauthorized, err)
				}
				return err
			}

			ctx := u.NewContext(r.Context())

			next.ServeHTTP(w, r.WithContext(ctx))

			return nil
		})
	}
}

func newAuthRequiredMid(next http.Handler) http.Handler {
	return errHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		usr, ok := domain.UserFromContext(r.Context())
		if !ok {
			return newWebErr(authRequiredErrCode, http.StatusUnauthorized, nil)
		}

		if !*usr.IsActive {
			return newWebErr(inactiveUserErrCode, http.StatusUnauthorized, nil)
		}

		next.ServeHTTP(w, r)
		return nil
	})
}

func newAdminOnlyMid(next http.Handler) http.Handler {
	return errHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		usr, ok := domain.UserFromContext(r.Context())
		if !ok {
			return errors.New("user can't get from request's context")
		}
		if !*usr.IsActive {
			return errors.New("inactive user")
		}

		if *usr.IsAdmin != true {
			return newWebErr(unknownErrCode, http.StatusUnauthorized, errors.New("permission denied"))
		}

		next.ServeHTTP(w, r)
		return nil
	})
}
