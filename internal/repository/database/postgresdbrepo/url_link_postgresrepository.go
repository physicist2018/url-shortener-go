package postgresdbrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"

	"github.com/physicist2018/url-shortener-go/internal/domain"
)

type PostgresDBLinkRepository struct {
	db *sql.DB
}

func NewDBLinkRepository(connStr string) (*PostgresDBLinkRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	dblink := &PostgresDBLinkRepository{db: db}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dblink.create(ctx); err != nil {
		return nil, err
	}

	return dblink, nil
}

func (d *PostgresDBLinkRepository) Store(ctx context.Context, urllink *domain.URLLink) error {
	query := `INSERT INTO links(short_url, original_url) VALUES($1, $2);`
	_, err := d.db.ExecContext(ctx, query, urllink.ShortURL, urllink.LongURL)
	return err
}

func (d *PostgresDBLinkRepository) Find(ctx context.Context, shortURL string) (*domain.URLLink, error) {
	query := `SELECT short_url, original_url FROM links WHERE short_url=$1;`
	row := d.db.QueryRowContext(ctx, query, shortURL)

	var urllink domain.URLLink
	err := row.Scan(&urllink.ShortURL, &urllink.LongURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("короткая сссылка не найдена")
		}
		return nil, err
	}
	return &urllink, nil
}

func (d *PostgresDBLinkRepository) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func (d *PostgresDBLinkRepository) create(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS links (
    id SERIAL PRIMARY KEY,
    short_url TEXT NOT NULL UNIQUE,
    original_url TEXT NOT NULL);`
	_, err := d.db.ExecContext(ctx, query)
	return err

}

func (d *PostgresDBLinkRepository) Close() error {
	return d.db.Close()
}
