package main

import (
	"context"
	"log"
	"makar/stemmer/pkg/database"
	"makar/stemmer/pkg/xkcd"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"makar/stemmer/pkg/config"
	"makar/stemmer/pkg/requests"
)

// setup функция для инициализации зависимостей
func setup() *requests.App {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		<-sigs
		cancelFunc()
	}()

	const configPath = "../../config_test.yaml"
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Init(cfg.DBFile)
	if err != nil {
		log.Fatal(err)
	}

	client, err := xkcd.Init(cfg.SourceURL)
	if err != nil {
		log.Fatal(err)
	}

	app := &requests.App{
		Db:     db,
		Client: client,
	}

	err = requests.DBDownloadComics(app, ctx, cfg.Parallel, cfg.IndexFile)
	if err != nil {
		log.Fatal(err)
	}

	return app
}

func BenchmarkFindComicsDB(b *testing.B) {
	b.ResetTimer()
	app := setup()
	input := "find home"

	for i := 0; i < b.N; i++ {
		err := requests.DBFindComics(app, input, false, 10) // Здесь предполагается ограничение поиска
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFindComicsIndex(b *testing.B) {
	b.ResetTimer()
	app := setup()
	input := "find home"
	for i := 0; i < b.N; i++ {
		err := requests.DBFindComics(app, input, true, 10) // Здесь предполагается ограничение поиска
		if err != nil {
			b.Fatal(err)
		}
	}
}
