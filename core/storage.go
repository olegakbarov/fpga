package core

import (
	"github.com/olegakbarov/io.confs.core/domain"
)

type (
	Sorting uint

	UserRepository interface {
		Add(*domain.User) error
		One(uint) (*domain.User, error)
		OneByEmail(string) (*domain.User, error)
		ExistsByEmail(string) (bool, error)
		Update(*domain.User) error
	}

	StorageFactory interface {
		NewUserRepository() UserRepository
	}
)

const (
	SortByIDDesc Sorting = iota
)
