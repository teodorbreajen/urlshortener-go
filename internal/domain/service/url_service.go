package service

import (
	"context"
	"errors"

	"urlshortener/internal/domain/model"
	"urlshortener/internal/port"
)

var (
	ErrURLNotFound  = errors.New("URL not found")
	ErrDuplicateURL = errors.New("URL already exists")
	ErrInvalidURL   = errors.New("invalid URL")
)

type urlService struct {
	repo port.URLRepository
}

func NewURLService(repo port.URLRepository) port.URLService {
	return &urlService{repo: repo}
}

func (s *urlService) CreateShortURL(ctx context.Context, originalURL string) (*model.URL, error) {
	if originalURL == "" {
		return nil, ErrInvalidURL
	}

	existing, err := s.repo.FindByOriginalURL(ctx, originalURL)
	if err == nil && existing != nil {
		return existing, nil
	}

	url, err := model.NewURL(originalURL)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, url); err != nil {
		return nil, err
	}

	return url, nil
}

func (s *urlService) GetByCode(ctx context.Context, code string) (*model.URL, error) {
	if code == "" {
		return nil, ErrInvalidURL
	}

	url, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if url == nil {
		return nil, ErrURLNotFound
	}

	return url, nil
}

func (s *urlService) GetByOriginalURL(ctx context.Context, originalURL string) (*model.URL, error) {
	if originalURL == "" {
		return nil, ErrInvalidURL
	}

	url, err := s.repo.FindByOriginalURL(ctx, originalURL)
	if err != nil {
		return nil, err
	}
	if url == nil {
		return nil, ErrURLNotFound
	}

	return url, nil
}

func (s *urlService) IncrementClicks(ctx context.Context, code string) (string, error) {
	url, err := s.GetByCode(ctx, code)
	if err != nil {
		return "", err
	}

	url.IncrementClicks()

	if err := s.repo.Update(ctx, url); err != nil {
		return "", err
	}

	return url.OriginalURL, nil
}

func (s *urlService) List(ctx context.Context, offset, limit int) ([]*model.URL, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.FindAll(ctx, offset, limit)
}
