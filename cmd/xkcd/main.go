package main

import (
	"context"
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

	// Считывание конфига
	const configPath = "config.yaml"
	config, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Создание БД, подгрузка конфига в пакеты
	initAll(config)

	err = requests.DBDownloadComics(app, ctx, config.Parallel)
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
