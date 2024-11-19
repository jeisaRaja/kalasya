package models

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrEmailDuplicate     = errors.New("user with this email exists")
	ErrSubdomainDuplicate = errors.New("blog with this subdomain exists")
	ErrInvalidCredentials = errors.New("email or password invalid")
)

type Models struct {
	Users interface {
		Insert(u *User) error
		Get(id int64) (*User, error)
		Exists(user *User) error
		Authenticate(email, password string) (int, error)
	}
	Blogs interface {
		Get(subdomain string) (*Blog, *BlogPost, error)
	}
	BlogPost interface {
		Update(blog *Blog, post *BlogPost) error
	}
}

func New(db *sql.DB) Models {
	return Models{
		Users:    UserModel{DB: db},
		Blogs:    BlogModel{DB: db},
		BlogPost: BlogPostModel{DB: db},
	}
}
