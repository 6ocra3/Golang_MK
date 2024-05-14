package json

import (
	"encoding/json"
	"os"
)

func BuildIndex(db *Database, indexFile string) error {
	indexData := makeIndex(db)
	db.Index = indexData
	err := saveIndex(indexData, indexFile)
	return err
}

func LoadIndex(db *Database, indexFile string) error {
	indexData, err := loadIndex(indexFile)
	if os.IsNotExist(err) {
		err := BuildIndex(db, indexFile)
		return err
	}
	if err != nil {
		return err
	}
	db.Index = indexData
	return nil
}

func makeIndex(db *Database) map[string][]int {
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

func saveIndex(indexData map[string][]int, indexFile string) error {
	// Сохранение данных в файл
	file, err := os.Create(indexFile)
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

func loadIndex(indexFile string) (map[string][]int, error) {
	// Подгрузка данных из уже созданного индекса
	file, err := os.Open(indexFile)
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
