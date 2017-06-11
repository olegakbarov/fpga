package postgres

import (
	"database/sql"

	"github.com/olegakbarov/io.confs.core/core"
	db "upper.io/db.v2"
	"upper.io/db.v2/lib/sqlbuilder"
)

type (
	tableName string

	repository struct {
		sess sqlbuilder.Database
	}
)

const (
	users = `users`
)

func handleErr(err error) error {
	if err == db.ErrNoMoreRows || err == sql.ErrNoRows {
		return core.ErrNoRows
	}
	return err
}
