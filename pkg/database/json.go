package database

import (
	"encoding/json"

	"os"
)

type Comics struct {
	ID       int      `json:"-"`
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

type Database struct {
	Entries  map[int]*Comics
	Index    map[string][]int
	FilePath string
}

func Init(FilePath string) (*Database, error) {

	// Инициализация конфига. Создание файла БД, если его нет

	file, err := os.OpenFile(FilePath, os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	newDB := &Database{
		Entries:  make(map[int]*Comics),
		FilePath: FilePath,
	}

	// Загрузка комиксов в базу данных
	err = GetComics(newDB)
	if err != nil {
		return nil, err
	}

	return newDB, nil

}

func AddComics(db *Database, comics []*Comics) error {
	// Добавление новых комиксов к имеющимся

	for _, comic := range comics {
		db.Entries[comic.ID] = comic
	}

	err := UpdateDB(db)

	if err != nil {
		return err
	}

	return nil

}

func UpdateDB(db *Database) error {
	// Запись в JSON
	updatedData, err := json.MarshalIndent(db.Entries, "", " ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(db.FilePath, updatedData, 0644); err != nil {
		return err
	}

	return nil
}

func GetComics(db *Database) error {

	// Считывание всех комиксов из JSON

	var existingComics map[int]*Comics

	data, err := os.ReadFile(db.FilePath)

	if err != nil {
		return err
	}

	// Если комиксов нет, то вернуть пустой map
	if len(data) == 0 {
		existingComics = make(map[int]*Comics)
	} else {
		err = json.Unmarshal(data, &existingComics)

		if err != nil {
			return err
		}
	}

	db.Entries = existingComics

	return nil
}
