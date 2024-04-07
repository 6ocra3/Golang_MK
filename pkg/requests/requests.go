package requests

import (
	"encoding/json"
	"fmt"
	"makar/stemmer/pkg/database"
	"makar/stemmer/pkg/words"
	"makar/stemmer/pkg/xkcd"
)

type App struct {
	Db     *database.Database
	Client *xkcd.Client
}

func DBDownloadComics(app *App, start int, end int) error {

	// Подгрузка данных с сайта
	comics, err := xkcd.DownloadComics(app.Client, start, end)
	if err != nil {
		return err
	}
	fmt.Print("Комиксы загружены\n")
	// Обработка полученных данных
	formatedComics, err := formatComics(comics)
	if err != nil {
		return nil
	}
	fmt.Print("Комиксы обработаны\n")

	// Добавление данных в БД
	err = database.AddComics(app.Db, formatedComics)
	if err != nil {
		return err
	}
	fmt.Print("Комиксы сохранены\n")

	return nil
}

func DBPrintComics(app *App, limit int) error {
	// Формирование списка комиксов, с id меньше limit
	limitedComics := make(map[int]*database.Comics)
	for i := 1; i <= limit; i++ {
		_, ok := app.Db.Entries[i]
		if !ok {
			break
		}
		limitedComics[i] = app.Db.Entries[i]
	}

	// Вывод JSON в консоль

	jsonData, err := json.MarshalIndent(limitedComics, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))

	return nil

}

func formatComics(comics []*xkcd.RawComic) ([]*database.Comics, error) {
	fmt.Print("Начало преобразования")
	// Преобразование комикса в нужный формат. Стеминг, фильтр полей
	formatedComics := make([]*database.Comics, 0, len(comics))
	for _, comic := range comics {
		description := comic.Transcript + comic.Alt

		stemmedDescription, err := words.StemmString(description)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		formatComic := database.Comics{ID: comic.ID, Url: comic.Url, Keywords: stemmedDescription}
		formatedComics = append(formatedComics, &formatComic)
	}
	return formatedComics, nil
}
