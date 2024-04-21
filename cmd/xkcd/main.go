package main

import (
	"context"
	"flag"
	"log"
	"makar/stemmer/pkg/config"
	"makar/stemmer/pkg/database"
	"makar/stemmer/pkg/requests"
	"makar/stemmer/pkg/xkcd"
	"os"
	"os/signal"
	"syscall"
)

var app *requests.App

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		<-sigs
		cancelFunc()
	}()

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
	initAll(config)

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

func initAll(config *config.Config) {

	db, err := database.Init(config.DBFile)
	if err != nil {
		log.Fatal(err)
	}

	client, err := xkcd.Init(config.SourceURL)
	if err != nil {
		log.Fatal(err)
	}

	app = &requests.App{
		Db:     db,
		Client: client,
	}

}
