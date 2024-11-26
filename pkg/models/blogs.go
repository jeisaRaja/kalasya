package models

import (
	"database/sql"
	"fmt"
	"html/template"
	"time"
)

type Blog struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	Name       string `json:"blog_name"`
	MainPostID int64
	Subdomain  string `json:"subdomain"`
	Nav        string
	NavHTML    template.HTML
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	AuthorName string
}

type BlogModel struct {
	DB *sql.DB
}

func (m *BlogModel) CreateBlog(b *Blog) error {
	err := m.DB.QueryRow(`
        INSERT INTO blogs (user_id, name, subdomain, nav, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id;`,
		b.UserID, b.Name, b.Subdomain, b.Nav, time.Now().UTC(), time.Now().UTC()).Scan(&b.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m BlogModel) GetID(subdomain string) (*int64, error) {
	var id int64
	stmt := `SELECT id FROM blogs WHERE subdomain = $1`
	err := m.DB.QueryRow(stmt, subdomain).Scan(&id)
	if err == ErrRecordNotFound {
		return nil, ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}

	return &id, nil
}

func (m BlogModel) Get(subdomain string) (*Blog, *Post, error) {
	var blog Blog
	var blogPost Post
	query := `SELECT id, name, subdomain, nav, user_id, main_post_id, updated_at FROM blogs WHERE subdomain = $1`
	err := m.DB.QueryRow(query, subdomain).Scan(&blog.ID, &blog.Name, &blog.Subdomain, &blog.Nav, &blog.UserID, &blog.MainPostID, &blog.UpdatedAt)
	if err != nil {
		return nil, nil, fmt.Errorf("while fetching blog by subdomain: %v", err)
	}
	query = `SELECT name FROM users WHERE id = $1`
	err = m.DB.QueryRow(query, blog.UserID).Scan(&blog.AuthorName)
	if err != nil {
		return nil, nil, fmt.Errorf("while fetching user name by user id: %v, user id is %d", err, blog.UserID)
	}
	query = `SELECT title, content, created_at, updated_at FROM blog_posts WHERE id = $1`
	err = m.DB.QueryRow(query, blog.MainPostID).Scan(&blogPost.Title, &blogPost.Content, &blog.CreatedAt, &blogPost.UpdatedAt)
	if err != nil {
		return nil, nil, fmt.Errorf("while fetching post by post id: %v", err)
	}
	return &blog, &blogPost, nil
}
