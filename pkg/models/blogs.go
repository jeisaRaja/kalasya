package models

import "database/sql"

type BlogsModel struct {
	DB *sql.DB
}

type Blog struct {
	ID        int64
	UserID    int64
	Title     string
	Subdomain string
}

func (m BlogsModel) Insert(blog *Blog) error {
	return nil
}

func (m BlogsModel) Get(id int64) (*Blog, error) {
	return nil, nil
}

func (m BlogsModel) Exists(subdomain string) (bool, error) {
	var count int
	stmt := `SELECT COUNT(*) FROM blogs WHERE subdomain = $1`
	err := m.DB.QueryRow(stmt, subdomain).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
