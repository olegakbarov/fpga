package web

type (
	errCode byte
	webErr  struct {
		Code     errCode
		HTTPCode int
		inner    error
	}
)

const (
	noRowsErrCode errCode = iota
	validationErrCode
	invalidTokenErrCode
	expiredTokenErrCode
	wrongCredErrCode
	authRequiredErrCode
	inactiveUserErrCode
	emailExistsErrCode
	badRequestErrCode
	unknownErrCode
)

func newWebErr(c errCode, sc int, inner error) *webErr {
	return &webErr{c, sc, inner}
}

func (e *webErr) Error() string {
	switch e.Code {
	default:
		if e.inner != nil {
			return e.inner.Error()
		}
		return "not implemented yet"
	case noRowsErrCode:
		return "not found"
	case invalidTokenErrCode:
		return "invalid token"
	case expiredTokenErrCode:
		return "expired token"
	case wrongCredErrCode:
		return "wrong credentials"
	case authRequiredErrCode:
		return "auth required"
	case inactiveUserErrCode:
		return "inactive user"
	case emailExistsErrCode:
		return "email address already exists"
	case badRequestErrCode:
		return "bad request"
	}
}
