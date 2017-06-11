package postgres

import (
	"sync"

	"github.com/olegakbarov/io.confs.core/core"
	"upper.io/db.v2/lib/sqlbuilder"
)

type (
	storageFactory struct {
		sess sqlbuilder.Database
	}
)

var (
	userRepositoryInstance core.UserRepository
	userRepositoryOnce     sync.Once
)

func NewStorage(session sqlbuilder.Database) core.StorageFactory {
	return &storageFactory{session}
}

func (sf *storageFactory) NewUserRepository() core.UserRepository {
	userRepositoryOnce.Do(func() {
		userRepositoryInstance = NewUserRepository(sf.sess)
	})
	return userRepositoryInstance
}

// func (sf *storageFactory) NewCatalogRepository() core.CatalogRepository {
//     catalogRepositoryOnce.Do(func() {
//         catalogRepositoryInstance = NewCatalogRepository(sf.sess)
//     })
//     return catalogRepositoryInstance
// }

// func (sf *storageFactory) NewImageRepository() core.ImageRepository {
//     imageRepositoryOnce.Do(func() {
//         imageRepositoryInstance = newImageRepository()
//     })
//     return imageRepositoryInstance
// }
