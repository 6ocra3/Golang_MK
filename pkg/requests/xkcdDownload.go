package requests

import (
	"context"
	"fmt"
	"makar/stemmer/pkg/database"
	"makar/stemmer/pkg/words"
	"makar/stemmer/pkg/xkcd"
	"sync"
)

type ParseError struct {
	id  int
	err error
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

func downloadComics(app *App, ctx context.Context, parallel int) *[]*database.Comics {
	errCnt := 0
	comics := make([]*database.Comics, 0)
	counter := nextId(app.Db)
	tasks := make(chan int)
	results := make(chan *database.Comics)
	errors := make(chan ParseError)
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
		case <-errors:
			//fmt.Printf("Err-%d ", err.id)
			errCnt++
			// Остановка подгрузки если все горутины закончатся с ошибкой
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

func worker(wg *sync.WaitGroup, app *App, tasks <-chan int, results chan<- *database.Comics, errors chan<- ParseError, done <-chan struct{}, ctx context.Context) {
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
				errors <- ParseError{id, err}
				wg.Done()
				return
			}

			formatComic, err := formatComics(comic)
			if err != nil {
				errors <- ParseError{id, err}
				wg.Done()
				return
			}
			select {
			case <-done:
				wg.Done()
				return
			case results <- formatComic:
			}
		}
	}
}

func formatComics(comic *xkcd.RawComic) (*database.Comics, error) {
	// Преобразование комикса в нужный формат. Стеминг, фильтр полей
	description := comic.Transcript + comic.Alt

	stemmedDescription, err := words.StemmString(description)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	formatComic := database.Comics{ID: comic.ID, Url: comic.Url, Keywords: stemmedDescription}
	return &formatComic, nil
}
