package core

import (
	"errors"
	"fmt"
)

type (
	ValidationErr struct {
		msg  string
		args []interface{}
	}

	TokenErr struct {
		msg     string
		expired bool
	}
)

var (
	ErrNoRows           = errors.New("not found")
	ErrEmailExists      = errors.New("email already exists")
	ErrEmailNotExists   = errors.New("email not exists")
	ErrInActiveUser     = errors.New("inactive user")
	ErrWrongCredentials = errors.New("wrong credentials")
)

func NewValidationErr(msg string, args ...interface{}) error {
	return &ValidationErr{msg, args}
}

func NewTokenErr(msg string, expired bool) error {
	return &TokenErr{msg, expired}
}

func (e *ValidationErr) Error() string {
	return fmt.Sprintf(e.msg, e.args...)
}

func (e *TokenErr) Error() string {
	return e.msg
}

func (e *TokenErr) Expired() bool {
	return e.expired
}
