package services

import "github.com/jeisaraja/kalasya/pkg/models"

func (s *Service) GetPosts(subdomain string) ([]*models.Post, error) {
	blogID, err := s.blogs.GetID(subdomain)
	if err != nil {
		return nil, err
	}
	posts, err := s.posts.GetPosts(*blogID)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
