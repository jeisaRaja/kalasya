package models

import (
	"database/sql"
	"net/mail"
	"time"

	"github.com/jeisaraja/kalasya/pkg/forms"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Name         string
	Email        string
	Password     string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRegistration struct {
	ID        int
	Subdomain string
	BlogName  string
	Name      string
	Email     string
	Password  string
}

type UserLogin struct {
	Email    string
	Password string
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) CreateUser(u *UserRegistration) error {
	err := m.Exists(u.Email, u.Subdomain)
	if err != nil {
		return err
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}
	_, err = m.DB.Exec(`
        INSERT INTO users (email, name, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5);`, u.Email, u.Name, passwordHash, time.Now().UTC(), time.Now().UTC())

	if err != nil {
		return err
	}

	return nil
}

func (m UserModel) Exists(email, subdomain string) error {
	var count int
	stmt := `SELECT COUNT(*) from users WHERE email = $1`
	err := m.DB.QueryRow(stmt, email).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrEmailDuplicate
	}
	stmt = `SELECT COUNT(*) from blogs WHERE subdomain = $1`
	err = m.DB.QueryRow(stmt, subdomain).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrSubdomainDuplicate
	}
	return nil
}

func (m UserModel) Get(email string) (*User, error) {
	var user User
	row := m.DB.QueryRow("SELECT id, name, email, password_hash FROM users WHERE email = $1", email)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, ErrInvalidCredentials
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func ValidateUserRegistration(form *forms.Form) {
	form.Required("blogname", "subdomain", "email", "password")
	form.MaxLength("blogname", 100)
	form.MaxLength("subdomain", 50)
	form.MinLength("blogname", 5)
	form.MinLength("subdomain", 3)
	form.MinLength("password", 8)
	form.MaxLength("password", 30)
	form.CheckFunc("email", EmailValid, "Email format invalid")
}

func EmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}
