package services

import (
	"time"

	"github.com/jeisaraja/kalasya/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) AuthenticateUser(email, password string) (*int, error) {
	u, err := s.users.Get(email)
	err = bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password))
	if err != nil {
		return nil, models.ErrInvalidCredentials
	}

	return &u.ID, nil
}

func (s *Service) CreateUserWithBlog(u *models.UserRegistration) error {
	tx, err := s.users.DB.Begin()
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

	err = s.users.CreateUser(u)
	if err != nil {
		return err
	}

	blog := models.Blog{
		UserID:    u.ID,
		Name:      u.BlogName,
		Subdomain: u.Subdomain,
		Nav:       "[Home](/) [Posts](/posts/)",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = s.blogs.CreateBlog(&blog)
	if err != nil {
		return err
	}

	post := models.Post{
		BlogID:    blog.ID,
		Title:     "Welcome to my blog!",
		Content:   "This is your first post.",
		Published: true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = s.posts.CreatePost(&post)
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
