package sqlite

import (
	"encoding/json"
	"makar/stemmer/pkg/database"
)

func (db *Database) LoadIndex() error {
	return nil
}

func (db *Database) BuildIndex() error {

	rows, err := db.db.Query("SELECT id, url, keywords FROM comics")
	if err != nil {
		return err
	}
	defer rows.Close()

	indexData := make(map[string][]int)

	for rows.Next() {
		var comic database.Comics
		var keywordsJSON []byte
		err := rows.Scan(&comic.ID, &comic.Url, &keywordsJSON)
		if err != nil {
			return err
		}

		err = json.Unmarshal(keywordsJSON, &comic.Keywords)
		if err != nil {
			return err
		}

		for _, word := range comic.Keywords {
			if indexData[word] == nil {
				indexData[word] = make([]int, 0)
			}
			indexData[word] = append(indexData[word], comic.ID)
		}

	}

	err = db.saveIndex(&indexData)

	return err
}

func (db *Database) saveIndex(indexData *map[string][]int) error {
	_, err := db.db.Exec("TRUNCATE TABLE index_table")
	if err != nil {
		return err
	}

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO index_table (word, ids) VALUES (?, ?)")
	if err != nil {

		return err
	}
	defer stmt.Close()

	for word, ids := range *indexData {
		idsJSON, err := json.Marshal(ids)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(word, idsJSON)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (db *Database) GetIds(word string) ([]int, error) {
	var idsJSON string
	err := db.db.QueryRow("SELECT ids FROM index_table WHERE word = ?", word).Scan(&idsJSON)
	if err != nil {
		return nil, err
	}

	var ids []int
	err = json.Unmarshal([]byte(idsJSON), &ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}
