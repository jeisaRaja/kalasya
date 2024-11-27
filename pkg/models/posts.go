package models

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"time"
)

type Post struct {
	ID        int
	BlogID    int
	Slug      string
	Title     string
	Content   string
	Published bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostView struct {
	Title     string
	Content   template.HTML
	Published bool
	CreatedAt time.Time
	UpdatedAt time.Time
	IsEdit    bool
}

type PostList struct {
	Title     string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostModel struct {
	DB *sql.DB
}

func (m *BlogModel) UpdateSelective(postID int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	query := "UPDATE blog_posts SET "
	args := []interface{}{}
	i := 1
	for column, value := range updates {
		if i > 1 {
			query += ", "
		}
		query += column + " = $" + strconv.Itoa(i)
		args = append(args, value)
		i++
	}
	curTime := time.Now()
	args = append(args, curTime)
	query += "updated_at = $" + strconv.Itoa(i)
	query += " WHERE id = $" + strconv.Itoa(i+1)
	args = append(args, postID)

	_, err := m.DB.Exec(query, args...)
	return err
}

func (m PostModel) Get() {

}

func (m PostModel) CreatePost(post *Post) error {
	stmt := `INSERT INTO blog_posts (slug, blog_id, title, content, published) VALUES ($1, $2, $3, $4, $5)`
	_, err := m.DB.Exec(stmt, post.Slug, post.BlogID, post.Title, post.Content, post.Published)
	if err != nil {
		return fmt.Errorf("failed to insert blog_posts: %v", err)
	}

	return nil
}

func (m PostModel) GetPostsBySubdomain(subdomain string) ([]*Post, error) {
	var posts []*Post
	stmt := `
      SELECT 
          bp.id, 
          bp.slug, 
          bp.title, 
          bp.created_at, 
          bp.updated_at
      FROM 
          blog_posts bp
      JOIN 
          blogs b ON bp.blog_id = b.id
      WHERE 
          b.subdomain = $1
          AND (bp.published = TRUE OR $2 = TRUE);
`
	rows, err := m.DB.Query(stmt, subdomain)
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
		err := rows.Scan(&post.ID, &post.Slug, &post.Title, &post.CreatedAt, &post.UpdatedAt)
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

func (m PostModel) GetByField(field string, value interface{}) (*Post, error) {
	query := fmt.Sprintf("SELECT id, blog_id, slug, title, content, published, created_at, updated_at FROM blog_posts WHERE %s = $1", field)
	var post Post
	err := m.DB.QueryRow(query, value).Scan(&post.ID, &post.BlogID, &post.Slug, &post.Title, &post.Content, &post.Published, &post.CreatedAt, &post.UpdatedAt)

	if err == ErrRecordNotFound {
		return nil, ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}

	return &post, nil
}
