package models

import (
	"database/sql"
	"fmt"
	"html/template"
	"time"
)

type BlogPost struct {
	ID            int64
	BlogID        int64
	Slug          string
	Title         string
	Content       string
	SanitizedHTML template.HTML
	Published     bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type BlogPostModel struct {
	DB *sql.DB
}

func (m BlogPostModel) Get() {

}

func (m BlogPostModel) Update(blog *Blog, post *BlogPost) error {
	stmt := `
    UPDATE blog_posts
    SET content = $1, updated_at = $2
    WHERE id = $3;`

	_, err := m.DB.Exec(stmt, post.Content, time.Now().UTC(), blog.MainPostID)
	if err != nil {
		return fmt.Errorf("failed to insert blog_posts: %v", err)
	}

	return nil
}
