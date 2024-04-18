package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"makar/stemmer/pkg/database"
	"makar/stemmer/pkg/words"
	"makar/stemmer/pkg/xkcd"
	"sync"
)

type MyError struct {
	id  int
	err error
}

type App struct {
	Db     *database.Database
	Client *xkcd.Client
}

func DBDownloadComics(app *App, ctx context.Context, parallel int) error {

	// Подгрузка данных с сайта
	comics := downloadComics(app, ctx, parallel)

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

func nextId(db *database.Database) func() int {
	cur := 1
	return func() int {
		for db.Entries[cur] != nil {
			cur++
		}
		cur++
		return cur - 1
	}
}

func downloadComics(app *App, ctx context.Context, parallel int) *[]*xkcd.RawComic {
	errCnt := 0
	comics := make([]*xkcd.RawComic, 0)
	counter := nextId(app.Db)
	tasks := make(chan int)
	results := make(chan *xkcd.RawComic)
	errors := make(chan MyError)
	done := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go worker(&wg, app, tasks, results, errors, done, ctx)
		tasks <- counter()

	}

downloadLoop:
	for {
		select {
		case <-ctx.Done():
			close(done)
			break downloadLoop
		case comic := <-results:
			comics = append(comics, comic)
			fmt.Printf("%d ", comic.ID)
			tasks <- counter()
		case err := <-errors:
			fmt.Printf("Err-%d ", err.id)
			errCnt++
			if errCnt == parallel {
				close(done)
				break downloadLoop
			}
		}

	}
	close(tasks)
	wg.Wait()
	close(results)
	close(errors)
	return &comics
}

func worker(wg *sync.WaitGroup, app *App, tasks <-chan int, results chan<- *xkcd.RawComic, errors chan<- MyError, done <-chan struct{}, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case <-done:
			wg.Done()
			return
		case id, ok := <-tasks:
			if !ok {
				return
			}
			comic, err := xkcd.DownloadComic(app.Client, id)
			if err != nil {
				errors <- MyError{id, err}
				continue
			}
			select {
			case <-done:
				wg.Done()
				return
			case results <- comic:
			}
		}
	}
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

func formatComics(comics *[]*xkcd.RawComic) ([]*database.Comics, error) {
	fmt.Print("Начало преобразования")
	// Преобразование комикса в нужный формат. Стеминг, фильтр полей
	formatedComics := make([]*database.Comics, 0, len(*comics))
	for _, comic := range *comics {
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
