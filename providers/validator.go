package providers

import (
	"net/mail"

	"github.com/alioygur/is"
	"github.com/olegakbarov/io.confs.core/core"
)

var newErr = core.NewValidationErr

func NewValidator() core.Validator {
	return &validator{}
}

type validator struct{}

func (v *validator) CheckEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return newErr("invalid email address")
	}
	return nil
}

func (v *validator) CheckRequired(val, field string) error {
	if len(val) != 0 {
		return nil
	}
	return newErr("the %s field is required", field)
}

func (v *validator) CheckStringLen(val string, min int, max int, field string) error {
	if is.StringLength(val, min, max) {
		return nil
	}
	return newErr("the %s field length must between %d and %d", field, min, max)
}
