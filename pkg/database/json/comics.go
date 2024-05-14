package json

import (
	"encoding/json"
	"makar/stemmer/pkg/database"

	"os"
)

type Database struct {
	Entries   map[int]*database.Comics
	Index     map[string][]int
	FilePath  string
	IndexFile string
}

func Init(FilePath string, IndexFile string) (*Database, error) {

	// Инициализация конфига. Создание файла БД, если его нет

	file, err := os.OpenFile(FilePath, os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	newDB := &Database{
		Entries:   make(map[int]*database.Comics),
		FilePath:  FilePath,
		IndexFile: IndexFile,
	}

	// Загрузка комиксов в базу данных
	comics, err := newDB.GetComics()
	if err != nil {
		return nil, err
	}

	newDB.Entries = comics

	return newDB, nil

}

func (db Database) AddComics(comics []*database.Comics) error {
	// Добавление новых комиксов к имеющимся

	for _, comic := range comics {
		db.Entries[comic.ID] = comic
	}

	err := db.UpdateDB()

	if err != nil {
		return err
	}

	return nil

}

func (db Database) UpdateDB() error {
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

func (db *Database) GetComics() (map[int]*database.Comics, error) {

	// Считывание всех комиксов из JSON

	var existingComics map[int]*database.Comics

	data, err := os.ReadFile(db.FilePath)

	if err != nil {
		return nil, err
	}

	// Если комиксов нет, то вернуть пустой map
	if len(data) == 0 {
		existingComics = make(map[int]*database.Comics)
	} else {
		err = json.Unmarshal(data, &existingComics)

		if err != nil {
			return nil, err
		}
	}

	return existingComics, nil
}

func (db Database) GetComic(id int) *database.Comics {
	return db.Entries[id]
}

func (db Database) CountComics() int {
	return len(db.Entries)
}
