package requests

import (
	"encoding/json"
	"fmt"
	"makar/stemmer/pkg/database"
	"makar/stemmer/pkg/xkcd"
)

func DBDownloadComics(start int, end int) error {

	// Подгрузка данных с сайта
	comics, err := xkcd.DownloadComics(start, end)
	if err != nil {
		return err
	}

	// Добавление данных в БД
	err = database.AddComics(comics)
	if err != nil {
		return err
	}

	return nil
}

func DBPrintComics(limit int) error {
	// Получение данных с БД
	comics, err := database.GetComics()
	if err != nil {
		return err
	}

	// Формирование списка комиксов, с id меньше limit
	limitedComics := make(map[int]*xkcd.Comic)

	if limit == -1 {
		limitedComics = comics
	} else {
		for i := 1; i <= limit; i++ {
			limitedComics[i] = comics[i]
		}
	}

	// Вывод JSON в консоль

	jsonData, err := json.MarshalIndent(limitedComics, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))

	return nil

}
