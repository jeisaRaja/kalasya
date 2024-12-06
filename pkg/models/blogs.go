package models

import (
	"database/sql"
	"html/template"
	"time"
)

type Blog struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Subdomain string    `json:"subdomain" db:"subdomain"`
	Nav       string    `json:"nav" db:"nav"`
	HomeID    int       `json:"main_post_id" db:"main_post_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type BlogView struct {
	Name      string
	Author    string
	Subdomain string
	NavHTML   template.HTML
	CreatedAt time.Time
	UpdatedAt time.Time
	Posts     PostList
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

func (m *BlogModel) GetBlogView(subdomain string) (*BlogView, error) {
	var blog BlogView
	err := m.DB.QueryRow(`
      SELECT 
            b.name AS BlogName,
            u.name AS AuthorName,
            b.subdomain AS Subdomain,
            b.nav AS NavHTML,
            b.created_at AS CreatedAt,
            b.updated_at AS UpdatedAt
        FROM 
            blogs b
        JOIN 
            users u ON b.user_id = u.id
        WHERE 
            b.subdomain = $1;
  `, subdomain).Scan(&blog.Name, &blog.Author, &blog.Subdomain, &blog.NavHTML, &blog.CreatedAt, &blog.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}

	return &blog, nil
}
