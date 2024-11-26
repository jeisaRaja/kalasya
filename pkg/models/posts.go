package models

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"time"
)

type Post struct {
	ID          int64
	BlogID      int64
	Slug        string
	Title       string
	Content     string
	Subdomain   string
	ContentHTML template.HTML
	Published   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsEdit      bool
}

type PostModel struct {
	DB *sql.DB
}

func (m PostModel) Get() {

}

func (m PostModel) Update(blog *Blog, post *Post) error {
	stmt := `
    UPDATE blog_posts
    SET content = $1, updated_at = $2
    WHERE id = $3;`

	_, err := m.DB.Exec(stmt, post.Content, time.Now().UTC(), blog.MainPostID)
	if err != nil {
		return fmt.Errorf("failed to update blog_posts: %v", err)
	}

	return nil
}

func (m PostModel) CreatePost(post *Post) error {
	stmt := `INSERT INTO blog_posts (slug, blog_id, title, content, published) VALUES ($1, $2, $3, $4, $5)`
	_, err := m.DB.Exec(stmt, post.Slug, post.BlogID, post.Title, post.Content, post.Published)
	if err != nil {
		return fmt.Errorf("failed to insert blog_posts: %v", err)
	}

	return nil
}

func (m PostModel) GetPosts(blogID int64) ([]*Post, error) {
	var posts []*Post
	stmt := `SELECT id, blog_id, slug, title 
         FROM blog_posts 
         WHERE blog_id = $1 AND slug <> '' AND title <> ''`
	rows, err := m.DB.Query(stmt, blogID)
	if err == sql.ErrNoRows {
		log.Printf("no rows found")
		return nil, err
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.BlogID, &post.Slug, &post.Title)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration error: %w", err)
	}

	return posts, nil
}

func (m PostModel) GetBySlug(slug string) (*Post, error) {
	var post Post
	stmt := `SELECT id, blog_id, slug, title, content, published, created_at, updated_at FROM blog_posts WHERE slug = $1`
	err := m.DB.QueryRow(stmt, slug).Scan(&post.ID, &post.BlogID, &post.Slug, &post.Title, &post.Content, &post.Published, &post.CreatedAt, &post.UpdatedAt)
	if err == ErrRecordNotFound {
		return nil, ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}
	return &post, nil
}
