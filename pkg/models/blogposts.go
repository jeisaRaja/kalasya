package models

import (
	"database/sql"
	"time"
)

type BlogPost struct {
	ID        int64
	BlogID    int64
	Slug      string
	Title     string
	Content   string
	Published bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BlogPostModel struct {
	DB *sql.DB
}

func (m BlogPostModel) Get(){

}
