package main

import (
	"context"
	"flag"
	"log"
	"makar/stemmer/pkg/config"
	"makar/stemmer/pkg/database"
	"makar/stemmer/pkg/requests"
	"makar/stemmer/pkg/xkcd"
	"os/signal"
	"syscall"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var input string
	var isIndexSearch bool
	flag.StringVar(&input, "s", "", "Флаг `-s` используется для ввода запроса")
	flag.BoolVar(&isIndexSearch, "i", false, "Флаг `-i` используется для включения поиска по индекс файлу")
	flag.Parse()

	// Считывание конфига
	const configPath = "config.yaml"
	config, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Создание БД, подгрузка конфига в пакеты
	app := initAll(config)

	// Подгрузка комиксов
	err = requests.DBDownloadComics(app, ctx, config.Parallel, config.IndexFile)
	if err != nil {
		log.Fatal(err)
	}

	err = requests.DBFindComics(app, input, isIndexSearch, config.SearchLimit)
	if err != nil {
		log.Fatal(err)
	}
}

func initAll(config *config.Config) *requests.App {

	db, err := database.Init(config.DBFile)
	if err != nil {
		log.Fatal(err)
	}

	client, err := xkcd.Init(config.SourceURL)
	if err != nil {
		log.Fatal(err)
	}

	app := &requests.App{
		Db:     db,
		Client: client,
	}

	return app

}
