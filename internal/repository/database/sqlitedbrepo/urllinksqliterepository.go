package sqlitedbrepo

import (
	"context"
	"database/sql"
	"errors"
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
		return nil, errors.Join(repoerrors.ErrorConnectingDB, err)
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Join(repoerrors.ErrorPingDB, err)
	}

	dblink := &SQLiteDBLinkRepository{db: db}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dblink.create(ctx); err != nil {
		return nil, err
	}

	return dblink, nil
}

func (d *SQLiteDBLinkRepository) Store(ctx context.Context, urllink domain.URLLink) (domain.URLLink, error) {
	queryInsert := `INSERT INTO links(user_id, short_url, original_url) VALUES($1, $2, $3);`
	_, err := d.db.ExecContext(ctx, queryInsert, urllink.UserID, urllink.ShortURL, urllink.LongURL)

	if err == nil {
		return urllink, nil
	}

	var sqliteError sqlite3.Error
	if !errors.As(err, &sqliteError) {
		return domain.URLLink{}, errors.Join(repoerrors.ErrorInsertShortLink, err)
	}

	// Обработка ошибки нарушения уникальности
	if sqliteError.ExtendedCode == sqlite3.ErrConstraintUnique {
		querySelect := `SELECT user_id, short_url, original_url FROM links WHERE original_url = $1 AND user_id = $2 LIMIT 1;`
		if err := d.db.GetContext(ctx, &urllink, querySelect, urllink.LongURL, urllink.UserID); err != nil {
			return domain.URLLink{}, errors.Join(repoerrors.ErrorSelectExistedShortLink, err)
		}
		return urllink, errors.Join(repoerrors.ErrorShortLinkAlreadyInDB, err)
	}

	// Обработка других ошибок SQLite
	return domain.URLLink{}, errors.Join(repoerrors.ErrorSQLInternal, err)
}

// change function specification
func (d *SQLiteDBLinkRepository) Find(ctx context.Context, shortURL, userID string) (domain.URLLink, error) {
	query := `SELECT user_id, short_url, original_url FROM links WHERE short_url=$1 AND user_id=$2 LIMIT 1;`
	var urllink domain.URLLink
	if err := d.db.GetContext(ctx, &urllink, query, shortURL, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// не найден короткий URL
			return domain.URLLink{}, errors.Join(repoerrors.ErrorShortLinkNotFound, err)
		}
		// при извлечении произошла ошибка
		return domain.URLLink{}, errors.Join(repoerrors.ErrorSelectExistedShortLink, err)
	}
	return urllink, nil
}

func (d *SQLiteDBLinkRepository) FindAll(ctx context.Context, userID string) ([]domain.URLLink, error) {
	query := `SELECT user_id, short_url, original_url FROM links WHERE user_id=$1;`
	var urllinks []domain.URLLink
	if err := d.db.SelectContext(ctx, &urllinks, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// не найден короткий URL
			return nil, errors.Join(repoerrors.ErrorShortLinkNotFound, err)
		}
		return nil, errors.Join(repoerrors.ErrorSelectShortLinks, err)
	}
	return urllinks, nil

}

func (d *SQLiteDBLinkRepository) Ping(ctx context.Context) error {
	if err := d.db.PingContext(ctx); err != nil {
		return errors.Join(repoerrors.ErrorPingDB, err)
	}
	return nil
}

func (d *SQLiteDBLinkRepository) create(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS links(
		ID INTEGER PRIMARY KEY,
		user_id TEXT NOT NULL,
		short_url TEXT NOT NULL,
		original_url TEXT NOT NULL UNIQUE
	);`
	if _, err := d.db.ExecContext(ctx, query); err != nil {
		return errors.Join(repoerrors.ErrorTableCreate, err)
	}
	return nil
}

func (d *SQLiteDBLinkRepository) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}
