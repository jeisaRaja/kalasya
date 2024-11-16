package models

import (
	"database/sql"
	"time"
)

type Blog struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	BlogName  string    `json:"blog_name"`
	Subdomain string    `json:"subdomain"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BlogModel struct {
	DB *sql.DB
}

func (m BlogModel) Get(subdomain string) (*Blog, error) {
	var blog Blog
	query := `SELECT name, subdomain, updated_at FROM blogs WHERE subdomain = $1`
	err := m.DB.QueryRow(query, subdomain).Scan(&blog.BlogName, &blog.Subdomain, &blog.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &blog, nil
}
