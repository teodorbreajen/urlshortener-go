package port

import (
	"context"

	"urlshortener/internal/domain/model"
)

type URLRepository interface {
	Save(ctx context.Context, url *model.URL) error
	FindByCode(ctx context.Context, code string) (*model.URL, error)
	FindByOriginalURL(ctx context.Context, originalURL string) (*model.URL, error)
	Update(ctx context.Context, url *model.URL) error
	FindAll(ctx context.Context, offset, limit int) ([]*model.URL, error)
}
