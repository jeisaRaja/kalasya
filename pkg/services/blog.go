package services

import (
	"fmt"

	"github.com/jeisaraja/kalasya/pkg/models"
)

func (s *Service) GetPosts(subdomain string, author bool) ([]*models.Post, error) {
	posts, err := s.posts.GetPostsBySubdomain(subdomain, author)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *Service) GetPost(slug string, author bool) (*models.PostView, error) {
	post, err := s.posts.Get(slug, author)
	if err == models.ErrRecordNotFound {
		return nil, fmt.Errorf("post not found: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve post: %w", err)
	}

	return post, nil
}

func (s *Service) GetBlogHome(subdomain string) (*models.PostView, error) {
	post, err := s.posts.GetHome(subdomain)
	if err != nil {
		return nil, err
	}

	return post, err
}

func (s *Service) GetBlogInfo(subdomain string) (*models.BlogView, error) {
	blog, err := s.blogs.GetBlogView(subdomain)
	if err != nil {
		return nil, err
	}
	return blog, nil
}

func (s *Service) UpdateBlogPost(postID int, updates map[string]interface{}) error {
  err := s.posts.UpdateSelective(postID, updates)
  return err
}
