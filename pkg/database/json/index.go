package json

import (
	"encoding/json"
	"os"
)

func (db Database) BuildIndex() error {
	indexData := db.makeIndex()
	db.Index = indexData
	err := db.saveIndex(indexData)
	return err
}

func (db Database) LoadIndex(indexFile string) error {
	indexData, err := db.loadIndex()
	if os.IsNotExist(err) {
		err := db.BuildIndex()
		return err
	}
	if err != nil {
		return err
	}
	db.Index = indexData
	return nil
}

func (db Database) makeIndex() map[string][]int {
	// Создание индекса
	indexData := make(map[string][]int)
	for id := range db.Entries {
		for _, word := range db.Entries[id].Keywords {
			if indexData[word] == nil {
				indexData[word] = make([]int, 0)
			}
			indexData[word] = append(indexData[word], id)
		}
	}
	return indexData
}

func (db Database) saveIndex(indexData map[string][]int) error {
	// Сохранение данных в файл
	file, err := os.Create(db.IndexFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // Set indent for pretty printing
	if err := encoder.Encode(indexData); err != nil {
		return err
	}
	return nil
}

func (db Database) loadIndex() (map[string][]int, error) {
	// Подгрузка данных из уже созданного индекса
	file, err := os.Open(db.IndexFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var indexData map[string][]int
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&indexData); err != nil {
		return nil, err
	}

	return indexData, nil
}

func (db Database) GetIds(word string) []int {
	return db.Index[word]
}
