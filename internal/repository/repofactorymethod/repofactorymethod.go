package repofactorymethod

import (
	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/repository/database/postgresdbrepo"
	"github.com/physicist2018/url-shortener-go/internal/repository/database/sqlitedbrepo"
	"github.com/physicist2018/url-shortener-go/internal/repository/inmemory"
)

type Repofactorymethod struct{}

func NewRepofactorymethod() *Repofactorymethod {
	return &Repofactorymethod{}
}
func (r *Repofactorymethod) createInMemoryRepo(dbname string) (*inmemory.InMemoryLinkRepository, error) {
	return inmemory.NewInMemoryLinkRepository(dbname)
}

func (r *Repofactorymethod) createSQLiteRepo(dbname string) (*sqlitedbrepo.SQLiteDBLinkRepository, error) {
	return sqlitedbrepo.NewDBLinkRepository(dbname)
}

func (r *Repofactorymethod) createPostgresRepo(connStr string) (*postgresdbrepo.PostgresDBLinkRepository, error) {
	return postgresdbrepo.NewDBLinkRepository(connStr)
}

// Фабричный метод для создания репозитория
func (r *Repofactorymethod) CreateRepo(repoType string, params string) (domain.URLLinkRepo, error) {
	switch repoType {
	case "inmemory":
		return r.createInMemoryRepo(params)
	case "sqlite":
		return r.createSQLiteRepo(params)
	case "postgres":
		return r.createPostgresRepo(params)
	default:
		return nil, nil
	}
}
