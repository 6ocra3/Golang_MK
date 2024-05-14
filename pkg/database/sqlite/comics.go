package sqlite

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"log"
	"makar/stemmer/pkg/database"
)

type Index struct {
	Word string `db:"word"`
	Ids  []int  `db:"ids"`
}

type Database struct {
	db *sqlx.DB
}

func InitSQLite(dsn1 string) (*Database, error) {

	databaseURL := "mysql://xkcd:xkcd@tcp(localhost:3306)/xkcd"
	migrationsDir := "file://./migrations"

	m, err := migrate.New(migrationsDir, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	dsn := "xkcd:xkcd@tcp(localhost:3306)/xkcd"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	database := &Database{db: db}
	return database, nil
}

func (db *Database) AddComics(comics []*database.Comics) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO comics (url, keywords) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, comic := range comics {
		keywordsJSON, err := json.Marshal(&comic.Keywords)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(comic.Url, keywordsJSON)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *Database) getComic(id int) *database.Comics {
	var comic database.Comics
	var keywordsJSON []byte

	// Извлечение данных из таблицы
	err := db.db.QueryRow("SELECT id, url, keywords FROM comics WHERE id = ?", id).Scan(&comic.ID, &comic.Url, &keywordsJSON)
	if err != nil {
		return nil
	}

	// Десериализация списка ключевых слов из JSON
	err = json.Unmarshal(keywordsJSON, &comic.Keywords)
	if err != nil {
		log.Fatalln(err)
	}

	return &comic
}
