package services

import "github.com/jeisaraja/kalasya/pkg/models"

func (s *Service) GetPosts(subdomain string) ([]*models.Post, error) {
	posts, err := s.posts.GetPostsBySubdomain(subdomain)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
