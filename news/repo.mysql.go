package news

import (
	"errors"
	"github.com/jmoiron/sqlx"
)


var (
	ErrNotFound = errors.New("not found")
)

type RepoMysql struct {
	DB *sqlx.DB
}

func (repo *RepoMysql) findAll(page, limit int) ([]News, int, error) {
	news := []News{}
	err := repo.DB.Select(&news, "SELECT * FROM news")
	if err != nil {
		return news, 0, ErrNotFound
	}

	length := len(news)

	first := page * limit
	if first < 0 {
		first = 0
	}

	last := first + limit
	if last > length {
		last = length
	}

	if first > length {
		return news, 0, ErrNotFound
	}

	return news[first:last], length, nil
}

func (repo *RepoMysql) findByAlias(alias string) (*News, error) {
	news := News{}
	err := repo.DB.Get(&news, "SELECT * FROM news WHERE alias=?", alias)
	if err != nil {
		return &news, ErrNotFound
	}

	return &news, nil
}