package sqlitedbrepo

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"

	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/repository/repoerrors"
)

//go:embed linktable.sql
var queryCreateTable string

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
	// Логика работы такая: пытаемся вставить короткую ссылку в БД, при этом если сокращаемый урл
	// уже там, мы возвращаем эту ссылку
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
		querySelect := `SELECT user_id, short_url, original_url FROM links WHERE original_url = $1 LIMIT 1;`
		if err := d.db.GetContext(ctx, &urllink, querySelect, urllink.LongURL); err != nil {

			return domain.URLLink{}, errors.Join(repoerrors.ErrorSelectExistedShortLink, err)
		}
		return urllink, errors.Join(repoerrors.ErrorShortLinkAlreadyInDB, err)
	}

	// Обработка других ошибок SQLite
	return domain.URLLink{}, errors.Join(repoerrors.ErrorSQLInternal, err)
}

// change function specification
func (d *SQLiteDBLinkRepository) Find(ctx context.Context, shortURL string) (domain.URLLink, error) {
	query := `SELECT user_id, short_url, original_url FROM links WHERE short_url=$1 LIMIT 1;`
	var urllink domain.URLLink
	if err := d.db.GetContext(ctx, &urllink, query, shortURL); err != nil {
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
	if _, err := d.db.ExecContext(ctx, queryCreateTable); err != nil {
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

func (d *SQLiteDBLinkRepository) MarkDeletedBatch(ctx context.Context, links []domain.URLLink) error {
	queryDelete := `
	UPDATE links
	SET is_deleted = true
	WHERE user_id = ? AND short_url = ?`

	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %v\n", err)
	}
	for _, link := range links {
		_, err = tx.ExecContext(ctx, queryDelete, link.UserID, link.ShortURL)

		if err != nil {
			log.Println("Ошибка при пометке на удаление")
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Ошибка при откате транзакции: %v\n", rbErr)
			}
			return err
		}
	}
	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		fmt.Printf("Ошибка при коммите транзакции: %v\n", err)
	}
	return err
}
