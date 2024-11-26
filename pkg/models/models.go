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
		GetUserPassword(email, password string) (*User, error)
	}
	Blogs interface {
		Get(subdomain string) (*Blog, *Post, error)
		GetID(subdomain string) (*int64, error)
	}
	Post interface {
		GetPosts(blogID int64) ([]*Post, error)
		GetBySlug(slug string) (*Post, error)
		Update(blog *Blog, post *Post) error
		CreatePost(post *Post) error
	}
}

func New(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db},
		Blogs: BlogModel{DB: db},
		Post:  PostModel{DB: db},
	}
}
