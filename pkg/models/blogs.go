package models

import "time"

type Blog struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	BlogName  string    `json:"blog_name"`
	Subdomain string    `json:"subdomain"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
