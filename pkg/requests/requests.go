package requests

import (
	"encoding/json"
	"fmt"
	"makar/stemmer/pkg/database"
	"makar/stemmer/pkg/xkcd"
)

func DBDownloadComics(start int, end int) error {
	comics, err := xkcd.DownloadComics(start, end)
	if err != nil {
		return err
	}

	err = database.AddComics(comics)
	if err != nil {
		return err
	}

	return nil
}

func DBPrintComics(limit int) error {
	comics, err := database.GetComics()
	if err != nil {
		return err
	}
	limitedComics := make(map[int]*xkcd.Comic)

	if limit == -1 {
		limitedComics = comics
	} else {
		for i := 1; i <= limit; i++ {
			limitedComics[i] = comics[i]
		}
	}

	jsonData, err := json.MarshalIndent(limitedComics, "", "    ")
	if err != nil {
		return err
	}

	// Вывод JSON в консоль
	fmt.Println(string(jsonData))

	return nil

}
