package port

import (
	"context"

	"urlshortener/internal/domain/model"
)

type URLService interface {
	CreateShortURL(ctx context.Context, originalURL string) (*model.URL, error)
	GetByCode(ctx context.Context, code string) (*model.URL, error)
	GetByOriginalURL(ctx context.Context, originalURL string) (*model.URL, error)
	IncrementClicks(ctx context.Context, code string) (string, error)
	List(ctx context.Context, offset, limit int) ([]*model.URL, error)
}
