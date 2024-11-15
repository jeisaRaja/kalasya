package postgres

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(blogtitle, subdomain, name, email, password string) error {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (email, password_hash) VALUES ($1,$2)`
	_, err = m.DB.Exec(stmt, email, passwordHash)
	if err != nil {
		return err
	}
	return nil
}
