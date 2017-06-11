package postgres

import (
	"time"

	"github.com/olegakbarov/io.confs.core/core"
	"github.com/olegakbarov/io.confs.core/domain"
	"upper.io/db.v2/lib/sqlbuilder"
)

type (
	userRepository struct {
		sess sqlbuilder.Database
	}
)

func NewUserRepository(sess sqlbuilder.Database) core.UserRepository {
	return &userRepository{sess: sess}
}

func (ur *userRepository) Add(u *domain.User) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	res, err := ur.sess.InsertInto(users).Values(u).Exec()
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	u.ID = uint(id)
	return err
}

func (ur *userRepository) One(id uint) (*domain.User, error) {
	var u domain.User
	return &u, handleErr(ur.sess.SelectFrom(users).Where(`id=?`, id).One(&u))
}

func (ur *userRepository) OneByEmail(email string) (*domain.User, error) {
	var u domain.User
	return &u, handleErr(ur.sess.SelectFrom(users).Where(`email=?`, email).One(&u))
}

func (ur *userRepository) ExistsByEmail(email string) (bool, error) {
	row, err := ur.sess.QueryRow(`SELECT COUNT(id) FROM users WHERE email=?`, email)
	if err != nil {
		return false, err
	}
	var n int
	err = row.Scan(&n)
	return n > 0, err
}

func (ur *userRepository) Update(u *domain.User) error {
	u.UpdatedAt = time.Now()
	_, err := ur.sess.Update(users).Set(u).Where(`id=?`, u.ID).Exec()
	return err
}

// func GetOne(id string) (User error) {
//     var rec users.RawUser

//     err := db.QueryRow("SELECT * FROM users WHERE id=$1 ORDER BY id", id).Scan(
//         &rec.FirstName,
//         &rec.LastName,
//         &rec.Email,
//         &rec.Locale,
//         &rec.City,
//         &rec.Userpic,
//         &rec.Settings,
//         &rec.PasswordHash,
//         &rec.Deleted,
//         &rec.Created_at,
//         &rec.Updated_at,
//     )

//     if err != nil {
//         log.Fatal("Error quering the db- " + err.Error())
//         w.WriteHeader(500)
//         return
//     }

//     return rec.PublicFields(), nil
// }
