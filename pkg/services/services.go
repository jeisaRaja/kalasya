package services

import (
	"database/sql"

	"github.com/jeisaraja/kalasya/pkg/models"
)

type Service struct {
	users models.UserModel
	blogs models.BlogModel
	posts models.PostModel
}

func New(db *sql.DB) *Service {
	return &Service{
		users: models.UserModel{DB: db},
		blogs: models.BlogModel{DB: db},
		posts: models.PostModel{DB: db},
	}
}
