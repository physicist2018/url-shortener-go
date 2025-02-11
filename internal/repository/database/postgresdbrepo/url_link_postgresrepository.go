package postgresdbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/repository/repoerrors"
)

type PostgresDBLinkRepository struct {
	db *sqlx.DB
}

func NewDBLinkRepository(connStr string) (*PostgresDBLinkRepository, error) {
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	dblink := &PostgresDBLinkRepository{db: db}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dblink.create(ctx); err != nil {
		return nil, fmt.Errorf("ошибка создания таблицы: %w", err)
	}

	return dblink, nil
}

func (d *PostgresDBLinkRepository) Store(ctx context.Context, urllink *domain.URLLink) error {
	query := `INSERT INTO links(short_url, original_url) VALUES($1, $2);`
	_, err := d.db.ExecContext(ctx, query, urllink.ShortURL, urllink.LongURL)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				query_select := `SELECT short_url FROM links WHERE original_url = $1 LIMIT 1;`
				row := d.db.QueryRowContext(ctx, query_select, urllink.LongURL)
				row.Scan(&urllink.ShortURL)
				return repoerrors.ErrUrlAlreadyInDB
			}
		}
		// оборачиваем исходную ошибку
		return fmt.Errorf("какая-то непредвиденная ошибка %w", err)
	}
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
    short_url TEXT NOT NULL,
    original_url TEXT NOT NULL UNIQUE);`
	_, err := d.db.ExecContext(ctx, query)
	return err

}

func (d *PostgresDBLinkRepository) Close() error {
	return d.db.Close()
}
