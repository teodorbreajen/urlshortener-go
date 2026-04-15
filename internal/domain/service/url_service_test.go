package service

import (
	"context"
	"testing"

	"urlshortener/internal/domain/model"
)

type mockURLRepository struct {
	urls map[string]*model.URL
}

func newMockRepo() *mockURLRepository {
	return &mockURLRepository{urls: make(map[string]*model.URL)}
}

func (m *mockURLRepository) Save(ctx context.Context, url *model.URL) error {
	m.urls[url.Code] = url
	return nil
}

func (m *mockURLRepository) FindByCode(ctx context.Context, code string) (*model.URL, error) {
	return m.urls[code], nil
}

func (m *mockURLRepository) FindByOriginalURL(ctx context.Context, originalURL string) (*model.URL, error) {
	for _, url := range m.urls {
		if url.OriginalURL == originalURL {
			return url, nil
		}
	}
	return nil, nil
}

func (m *mockURLRepository) Update(ctx context.Context, url *model.URL) error {
	m.urls[url.Code] = url
	return nil
}

func (m *mockURLRepository) FindAll(ctx context.Context, offset, limit int) ([]*model.URL, error) {
	var result []*model.URL
	i := 0
	for _, url := range m.urls {
		if i >= offset && i < offset+limit {
			result = append(result, url)
		}
		i++
	}
	return result, nil
}

func TestCreateShortURL(t *testing.T) {
	repo := newMockRepo()
	svc := NewURLService(repo)

	url, err := svc.CreateShortURL(context.Background(), "https://example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if url == nil {
		t.Fatal("expected URL, got nil")
	}

	if url.Code == "" {
		t.Error("expected code, got empty")
	}

	if url.OriginalURL != "https://example.com" {
		t.Errorf("expected https://example.com, got %s", url.OriginalURL)
	}

	if url.Clicks != 0 {
		t.Errorf("expected 0 clicks, got %d", url.Clicks)
	}
}

func TestCreateShortURL_Duplicate(t *testing.T) {
	repo := newMockRepo()
	svc := NewURLService(repo)

	original, _ := svc.CreateShortURL(context.Background(), "https://example.com")

	duplicate, err := svc.CreateShortURL(context.Background(), "https://example.com")
	if err != nil {
		t.Fatalf("expected no error for duplicate, got %v", err)
	}

	if original.Code != duplicate.Code {
		t.Error("expected same code for duplicate URL")
	}
}

func TestCreateShortURL_InvalidURL(t *testing.T) {
	repo := newMockRepo()
	svc := NewURLService(repo)

	_, err := svc.CreateShortURL(context.Background(), "")
	if err != ErrInvalidURL {
		t.Errorf("expected ErrInvalidURL, got %v", err)
	}
}

func TestGetByCode(t *testing.T) {
	repo := newMockRepo()
	svc := NewURLService(repo)

	created, _ := svc.CreateShortURL(context.Background(), "https://example.com")

	found, err := svc.GetByCode(context.Background(), created.Code)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.Code != created.Code {
		t.Errorf("expected code %s, got %s", created.Code, found.Code)
	}
}

func TestGetByCode_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := NewURLService(repo)

	_, err := svc.GetByCode(context.Background(), "nonexistent")
	if err != ErrURLNotFound {
		t.Errorf("expected ErrURLNotFound, got %v", err)
	}
}

func TestIncrementClicks(t *testing.T) {
	repo := newMockRepo()
	svc := NewURLService(repo)

	url, _ := svc.CreateShortURL(context.Background(), "https://example.com")
	initialClicks := url.Clicks

	_, err := svc.IncrementClicks(context.Background(), url.Code)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, _ := svc.GetByCode(context.Background(), url.Code)
	if updated.Clicks != initialClicks+1 {
		t.Errorf("expected %d clicks, got %d", initialClicks+1, updated.Clicks)
	}
}

func TestList(t *testing.T) {
	repo := newMockRepo()
	svc := NewURLService(repo)

	for i := 0; i < 5; i++ {
		svc.CreateShortURL(context.Background(), "https://example.com/"+string(rune('a'+i)))
	}

	urls, err := svc.List(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(urls) != 5 {
		t.Errorf("expected 5 URLs, got %d", len(urls))
	}
}
