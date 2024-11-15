package models

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicate      = errors.New("record exists")
)

type Models struct {
	Users interface {
		Insert(u *User) error
		Get(id int64) (*User, error)
		Exists(email string) (bool, error)
	}
	Blogs interface {
		Insert(b *Blog) error
		Get(id int64) (*Blog, error)
		Exists(subdomain string) (bool, error)
	}
}

func New(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db},
		Blogs: BlogsModel{DB: db},
	}
}
