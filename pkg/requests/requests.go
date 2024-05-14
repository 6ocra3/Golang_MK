package requests

import (
	"context"
	"fmt"
	"makar/stemmer/pkg/database/json"
	"makar/stemmer/pkg/words"
	"makar/stemmer/pkg/xkcd"
)

type App struct {
	Db     *database.Database
	Client *xkcd.Client
}

func DBDownloadComics(app *App, ctx context.Context, parallel int, indexFile string) error {

	// Подгрузка данных с сайта
	comics := downloadComics(ctx, app, parallel)
	fmt.Print("\nКомиксы обработаны\n")

	// Добавление данных в БД
	err := database.AddComics(app.Db, comics)
	if err != nil {
		return err
	}
	fmt.Print("Комиксы сохранены\n")

	err = database.LoadIndex(app.Db, indexFile)
	if err != nil {
		return err
	}
	fmt.Print("Индекс построен\n")

	return nil
}

func DBFindComics(app *App, request string, limit int) (error, []string) {
	// Обрабатываем запрос
	stemRequest, err := words.StemmString(request)

	if err != nil {
		return err, nil
	}

	// Получение map keyword -> [id1, id2, id3]
	var searchResult map[string][]int
	searchResult = FindWithIndex(app, stemRequest)

	// Обработка map keyword -> [id1, id2, id3] и получение итогового списка id
	processedResult := processResult(app, searchResult, limit)

	// Выводим ссылки
	if len(processedResult) == 0 {
		return nil, nil
	}

	findedComics := make([]string, len(processedResult))

	for i := range processedResult {
		findedComics[i] = app.Db.Entries[processedResult[i]].Url
	}

	return nil, findedComics
}
