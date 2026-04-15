package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidURL = errors.New("invalid URL")
	ErrEmptyCode  = errors.New("code cannot be empty")
)

type URL struct {
	ID          string
	Code        string
	OriginalURL string
	Clicks      int64
	CreatedAt   time.Time
	ExpiresAt   *time.Time
}

func NewURL(originalURL string) (*URL, error) {
	if originalURL == "" {
		return nil, ErrInvalidURL
	}

	code := GenerateCode()

	return &URL{
		ID:          uuid.New().String(),
		Code:        code,
		OriginalURL: originalURL,
		Clicks:      0,
		CreatedAt:   time.Now(),
		ExpiresAt:   nil,
	}, nil
}

func (u *URL) IncrementClicks() {
	u.Clicks++
}

func (u *URL) IsExpired() bool {
	if u.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*u.ExpiresAt)
}

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateCode() string {
	id := uuid.New().String()
	code := make([]byte, 8)

	for i := 0; i < 8; i++ {
		idx := int(id[i]) % 62
		code[i] = base62Chars[idx]
	}

	return string(code)
}
