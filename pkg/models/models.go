package models

import (
	"database/sql"
	"errors"
)

const (
	Title     = "title"
	Content   = "content"
	Published = "published"
	HomeID    = "main_post_id"
	Name      = "name"
	Nav       = "nav"
	CreatedAt = "created_at"
	UpdatedAt = "updated_at"
)

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrEmailDuplicate     = errors.New("user with this email exists")
	ErrSubdomainDuplicate = errors.New("blog with this subdomain exists")
	ErrInvalidCredentials = errors.New("email or password invalid")
)

type Models struct {
	User UserModel
	Blog BlogModel
	Post PostModel
}

func New(db *sql.DB) Models {
	return Models{
		User: UserModel{DB: db},
		Blog: BlogModel{DB: db},
		Post: PostModel{DB: db},
	}
}
