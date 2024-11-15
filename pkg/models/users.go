package models

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string
	Password string
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(u *User) error {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (email, password_hash) VALUES ($1,$2)`
	_, err = m.DB.Exec(stmt, u.Email, passwordHash)
	if err != nil {
		return err
	}
	return nil
}

func (m UserModel) Get(id int64) (*User, error) {
	return nil, nil
}

func (m UserModel) Exists(email string) (bool, error) {
	var count int
	stmt := `SELECT COUNT(*) from users WHERE email = $1`
	err := m.DB.QueryRow(stmt, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
