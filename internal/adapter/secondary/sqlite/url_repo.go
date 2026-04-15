package sqlite

import (
	"context"
	"database/sql"
	"time"

	"urlshortener/internal/domain/model"
	"urlshortener/internal/port"

	_ "github.com/mattn/go-sqlite3"
)

type urlRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) port.URLRepository {
	return &urlRepository{db: db}
}

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createSchema(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS urls (
		id TEXT PRIMARY KEY,
		code TEXT UNIQUE NOT NULL,
		original_url TEXT NOT NULL,
		clicks INTEGER DEFAULT 0,
		created_at TEXT NOT NULL,
		expires_at TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_urls_code ON urls(code);
	CREATE INDEX IF NOT EXISTS idx_urls_original_url ON urls(original_url);
	`

	_, err := db.Exec(schema)
	return err
}

func (r *urlRepository) Save(ctx context.Context, url *model.URL) error {
	query := `INSERT INTO urls (id, code, original_url, clicks, created_at, expires_at) VALUES (?, ?, ?, ?, ?, ?)`

	var expiresAt *string
	if url.ExpiresAt != nil {
		t := url.ExpiresAt.Format(time.RFC3339)
		expiresAt = &t
	}

	_, err := r.db.ExecContext(ctx, query,
		url.ID,
		url.Code,
		url.OriginalURL,
		url.Clicks,
		url.CreatedAt.Format(time.RFC3339),
		expiresAt,
	)

	return err
}

func (r *urlRepository) FindByCode(ctx context.Context, code string) (*model.URL, error) {
	query := `SELECT id, code, original_url, clicks, created_at, expires_at FROM urls WHERE code = ?`

	var url model.URL
	var expiresAt sql.NullString

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&url.ID,
		&url.Code,
		&url.OriginalURL,
		&url.Clicks,
		&url.CreatedAt,
		&expiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if expiresAt.Valid {
		t, _ := time.Parse(time.RFC3339, expiresAt.String)
		url.ExpiresAt = &t
	}

	return &url, nil
}

func (r *urlRepository) FindByOriginalURL(ctx context.Context, originalURL string) (*model.URL, error) {
	query := `SELECT id, code, original_url, clicks, created_at, expires_at FROM urls WHERE original_url = ?`

	var url model.URL
	var expiresAt sql.NullString

	err := r.db.QueryRowContext(ctx, query, originalURL).Scan(
		&url.ID,
		&url.Code,
		&url.OriginalURL,
		&url.Clicks,
		&url.CreatedAt,
		&expiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if expiresAt.Valid {
		t, _ := time.Parse(time.RFC3339, expiresAt.String)
		url.ExpiresAt = &t
	}

	return &url, nil
}

func (r *urlRepository) Update(ctx context.Context, url *model.URL) error {
	query := `UPDATE urls SET clicks = ?, expires_at = ? WHERE code = ?`

	var expiresAt *string
	if url.ExpiresAt != nil {
		t := url.ExpiresAt.Format(time.RFC3339)
		expiresAt = &t
	}

	_, err := r.db.ExecContext(ctx, query, url.Clicks, expiresAt, url.Code)
	return err
}

func (r *urlRepository) FindAll(ctx context.Context, offset, limit int) ([]*model.URL, error) {
	query := `SELECT id, code, original_url, clicks, created_at, expires_at FROM urls ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []*model.URL
	for rows.Next() {
		var url model.URL
		var expiresAt sql.NullString

		if err := rows.Scan(
			&url.ID,
			&url.Code,
			&url.OriginalURL,
			&url.Clicks,
			&url.CreatedAt,
			&expiresAt,
		); err != nil {
			return nil, err
		}

		if expiresAt.Valid {
			t, _ := time.Parse(time.RFC3339, expiresAt.String)
			url.ExpiresAt = &t
		}

		urls = append(urls, &url)
	}

	return urls, nil
}
