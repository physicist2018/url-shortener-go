package postgresdbrepo

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	_ "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/repository/repoerrors"
)

//go:embed linktable.sql
var queryCreateTable string

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
func (d *PostgresDBLinkRepository) Store(ctx context.Context, urllink domain.URLLink) (domain.URLLink, error) {
	query := `INSERT INTO links(user_id, short_url, original_url) VALUES($1, $2, $3);`
	_, err := d.db.ExecContext(ctx, query, urllink.UserID, urllink.ShortURL, urllink.LongURL)

	if err == nil {
		return urllink, nil
	}

	var pqError *pq.Error
	if !errors.As(err, &pqError) {
		return domain.URLLink{}, errors.Join(repoerrors.ErrorInsertShortLink, err)
	}

	if pqError.Code == "23505" {
		querySelect := `SELECT user_id, short_url, original_url FROM links WHERE original_url = $1 LIMIT 1;`
		if err := d.db.GetContext(ctx, &urllink, querySelect, urllink.LongURL); err != nil {
			return domain.URLLink{}, errors.Join(repoerrors.ErrorSelectExistedShortLink, err)
		}
		return urllink, errors.Join(repoerrors.ErrorShortLinkAlreadyInDB, err)
	}

	// Обработка других ошибок Postgres
	return domain.URLLink{}, errors.Join(repoerrors.ErrorSQLInternal, err)
}

// TODO change function input parameters
func (d *PostgresDBLinkRepository) Find(ctx context.Context, shortURL string) (domain.URLLink, error) {
	query := `SELECT user_id, short_url, original_url, is_deleted FROM links WHERE short_url=$1 LIMIT 1;`
	var urllink domain.URLLink
	if err := d.db.GetContext(ctx, &urllink, query, shortURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// не найден короткий URL
			return domain.URLLink{}, errors.Join(repoerrors.ErrorShortLinkNotFound, err)
		}
		// при извлечении произошла ошибка
		return urllink, errors.Join(repoerrors.ErrorSelectExistedShortLink, err)
	}
	return urllink, nil
}

func (d *PostgresDBLinkRepository) FindAll(ctx context.Context, userID string) ([]domain.URLLink, error) {
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

func (d *PostgresDBLinkRepository) Ping(ctx context.Context) error {
	if err := d.db.PingContext(ctx); err != nil {
		return errors.Join(repoerrors.ErrorPingDB, err)
	}
	return nil
}

func (d *PostgresDBLinkRepository) create(ctx context.Context) error {
	if _, err := d.db.ExecContext(ctx, queryCreateTable); err != nil {
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

func (d *PostgresDBLinkRepository) MarkDeletedBatch(ctx context.Context, links domain.DeleteRecordTask) error {
	queryDelete := `UPDATE links SET is_deleted = true 	WHERE user_id = $1 AND short_url = ANY($2);`
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции : %v", err)
	}
	_, err = tx.ExecContext(ctx, queryDelete, links.UserID, pq.Array(links.ShortURLs))
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("ошибка отката транзакции : %v", rbErr)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка коммита транзакции : %v", err)
	}
	return nil
}
