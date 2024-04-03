package database

import (
	"encoding/json"
	"makar/stemmer/pkg/config"
	"makar/stemmer/pkg/xkcd"

	"os"
)

var configData *config.Config

func Init(config *config.Config) error {
	configData = config

	file, err := os.OpenFile(configData.DBFile, os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	file.Close()

	return nil

}

func AddComics(comics map[int]*xkcd.Comic) error {
	existingComics, err := GetComics()

	if err != nil {
		return err
	}

	for id, comic := range comics {
		existingComics[id] = comic
	}

	updatedData, err := json.MarshalIndent(existingComics, "", " ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configData.DBFile, updatedData, 0644); err != nil {
		return err
	}

	return nil

}

func GetComics() (map[int]*xkcd.Comic, error) {

	var existingComics map[int]*xkcd.Comic

	data, err := os.ReadFile(configData.DBFile)

	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		existingComics = make(map[int]*xkcd.Comic)
	} else {
		err = json.Unmarshal(data, &existingComics)

		if err != nil {
			return nil, err
		}
	}

	return existingComics, nil
}
