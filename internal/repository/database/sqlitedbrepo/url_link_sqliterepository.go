package sqlitedbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"

	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/repository/repoerrors"
)

type SQLiteDBLinkRepository struct {
	db *sqlx.DB
}

func NewDBLinkRepository(connStr string) (*SQLiteDBLinkRepository, error) {
	db, err := sqlx.Open("sqlite3", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	dblink := &SQLiteDBLinkRepository{db: db}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dblink.create(ctx); err != nil {
		return nil, err
	}

	return dblink, nil
}

func (d *SQLiteDBLinkRepository) Store(ctx context.Context, urllink *domain.URLLink) error {
	query := `INSERT INTO links(short_url, original_url) VALUES($1, $2);`
	_, err := d.db.ExecContext(ctx, query, urllink.ShortURL, urllink.LongURL)

	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			if sqliteError.ExtendedCode == sqlite3.ErrConstraintUnique {
				querySelect := `SELECT short_url FROM links WHERE original_url = $1 LIMIT 1;`
				row := d.db.QueryRowContext(ctx, querySelect, urllink.LongURL)
				row.Scan(&urllink.ShortURL)
				return repoerrors.ErrUrlAlreadyInDB
			}
		}
		// оборачиваем исходную ошибку
		return fmt.Errorf("какая-то непредвиденная ошибка %w", err)
	}
	return err
}

func (d *SQLiteDBLinkRepository) Find(ctx context.Context, shortURL string) (*domain.URLLink, error) {
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

func (d *SQLiteDBLinkRepository) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func (d *SQLiteDBLinkRepository) create(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS links(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short_url TEXT NOT NULL,
		original_url TEXT NOT NULL UNIQUE
	);`
	_, err := d.db.ExecContext(ctx, query)
	return err

}

func (d *SQLiteDBLinkRepository) Close() error {
	return d.db.Close()
}
