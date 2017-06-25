package news

import (
	"errors"

	"github.com/jmoiron/sqlx"
)


var (
	ErrNotFound = errors.New("not found")
	NoInserterFound = errors.New("No inserter found")
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

func (repo *RepoMysql) create(data News) bool {
	result := false

	tx, err := repo.DB.Begin()

	_, err = tx.Exec(
		"INSERT INTO news (title, alias, description, content) VALUES (?, ?, ?, ?)",
		data.Title, data.Alias, data.Description, data.Content,
	)

	if err != nil {
		result = false
	}

	tx.Commit()
	result = true

	return result
}

func (repo *RepoMysql) update(data News, alias string) {
	tx, _ := repo.DB.Begin()

	_, _ = tx.Exec(
		"UPDATE news SET title=?, alias=?, description=?, content=? WHERE alias=? LIMIT 1",
		data.Title, data.Alias, data.Description, data.Content, alias,
	)

	tx.Commit()
}
