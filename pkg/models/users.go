package models

import (
	"database/sql"
	"net/mail"
	"time"

	"github.com/jeisaraja/kalasya/pkg/forms"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	BlogName  string
	Subdomain string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(u *User) error {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}

	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.QueryRow(`
        INSERT INTO users (email, name, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		u.Email, u.Name, passwordHash, time.Now().UTC(), time.Now().UTC()).Scan(&u.ID)

	if err != nil {
		return err
	}

	nav := "[Home](/) [Posts](/posts/)"

	blog := Blog{
		UserID:    u.ID,
		Subdomain: u.Subdomain,
		Name:      u.BlogName,
		Nav:       nav,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = tx.QueryRow(`
        INSERT INTO blogs (user_id, name, subdomain, nav, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id;`,
		blog.UserID, blog.Name, blog.Subdomain, blog.Nav, blog.CreatedAt, blog.UpdatedAt).Scan(&blog.ID)

	if err != nil {
		return err
	}

	var post BlogPost
	err = tx.QueryRow(`
	INSERT INTO blog_posts (blog_id, title, content, published, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id;`,
		blog.ID, post.Title, post.Content, true, post.CreatedAt, post.UpdatedAt).Scan(&post.ID)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	UPDATE blogs
	SET main_post_id = $1
	WHERE id = $2;`,
		post.ID, blog.ID)

	if err != nil {
		return err
	}

	return nil
}

func (m UserModel) Get(id int64) (*User, error) {
	u := &User{}
	stmt := `SELECT id, name, email FROM users WHERE id = $1`
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email)
	if err == sql.ErrNoRows {
		return nil, ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}
	stmt = `SELECT name, subdomain FROM blogs WHERE user_id = $1`
	err = m.DB.QueryRow(stmt, id).Scan(&u.BlogName, &u.Subdomain)
	return u, nil
}

func (m UserModel) Exists(user *User) error {
	var count int
	stmt := `SELECT COUNT(*) from users WHERE email = $1`
	err := m.DB.QueryRow(stmt, user.Email).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrEmailDuplicate
	}
	stmt = `SELECT COUNT(*) from blogs WHERE subdomain = $1`
	err = m.DB.QueryRow(stmt, user.Subdomain).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrSubdomainDuplicate
	}
	return nil
}

func (m UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var passwordHash []byte
	row := m.DB.QueryRow("SELECT id, password_hash FROM users WHERE email = $1", email)
	err := row.Scan(&id, &passwordHash)
	if err == sql.ErrNoRows {
		return 0, ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	return id, nil
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
