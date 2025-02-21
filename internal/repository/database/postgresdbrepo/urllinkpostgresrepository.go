package postgresdbrepo

import (
	"context"
	"database/sql"
	"errors"
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
		return nil, errors.Join(repoerrors.ErrorConnectingDB, err)
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Join(repoerrors.ErrorPingDB, err)
	}

	dblink := &PostgresDBLinkRepository{db: db}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dblink.create(ctx); err != nil {
		return nil, err
	}

	return dblink, nil
}

// Store is a function that stores a URL link in the database.
// It takes a context and a URL link as arguments and returns an error.
// It inserts the URL link into the database using a prepared SQL query.
// If there is an error during the process, it checks if the error is a PostgreSQL error and if it is a unique constraint violation.
// If it is a unique constraint violation, it retrieves the short URL for the original URL from the database and returns a custom error.
// If there is any other error, it returns a formatted error with the original error.
func (d *PostgresDBLinkRepository) Store(ctx context.Context, urllink *domain.URLLink) error {
	query := `INSERT INTO links(short_url, original_url) VALUES($1, $2);`
	_, err := d.db.ExecContext(ctx, query, urllink.ShortURL, urllink.LongURL)

	if err == nil {
		return nil
	}

	var pqError *pq.Error
	if !errors.As(err, &pqError) {
		return errors.Join(repoerrors.ErrorInsertShortLink, err)
	}

	if pqError.Code == "23505" {
		querySelect := `SELECT short_url, original_url FROM links WHERE original_url = $1 LIMIT 1;`
		if err := d.db.GetContext(ctx, urllink, querySelect, urllink.LongURL); err != nil {
			return errors.Join(repoerrors.ErrorSelectExistedShortLink, err)
		}
		return errors.Join(repoerrors.ErrorShortLinkAlreadyInDB, err)
	}

	// Обработка других ошибок Postgres
	return errors.Join(repoerrors.ErrorSQLInternal, err)
}

func (d *PostgresDBLinkRepository) Find(ctx context.Context, shortURL string) (*domain.URLLink, error) {
	query := `SELECT short_url, original_url FROM links WHERE short_url=$1 LIMIT 1;`
	var urllink domain.URLLink
	if err := d.db.GetContext(ctx, &urllink, query, shortURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// не найден короткий URL
			return nil, errors.Join(repoerrors.ErrorShortLinkNotFound, err)
		}
		// при извлечении произошла ошибка
		return nil, errors.Join(repoerrors.ErrorSelectExistedShortLink, err)
	}
	return &urllink, nil
}

func (d *PostgresDBLinkRepository) Ping(ctx context.Context) error {
	if err := d.db.PingContext(ctx); err != nil {
		return errors.Join(repoerrors.ErrorPingDB, err)
	}
	return nil
}

func (d *PostgresDBLinkRepository) create(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS links (
        id SERIAL PRIMARY KEY,
        short_url TEXT NOT NULL,
        original_url TEXT NOT NULL UNIQUE);`
	if _, err := d.db.ExecContext(ctx, query); err != nil {
		return errors.Join(repoerrors.ErrorTableCreate, err)
	}
	return nil
}

func (d *PostgresDBLinkRepository) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}
