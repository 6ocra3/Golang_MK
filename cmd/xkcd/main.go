package main

import (
	"flag"
	"log"
	"makar/stemmer/pkg/config"
	"makar/stemmer/pkg/database"
	"makar/stemmer/pkg/requests"
	"makar/stemmer/pkg/xkcd"
	"math"
)

const base_batch_size = 10

var app *requests.App

func main() {

	// Считывание флагов
	printFlag := flag.Bool("o", false, "Флаг для вывода скаченных комиксов")
	limitFlag := flag.Int("n", math.MaxInt, "Флаг для установки лимита показываемых комиксов")

	flag.Parse()

	// Считывание конфига
	const configPath = "config.yaml"
	config, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Создание БД, подгрузка конфига в пакеты
	initAll(config)
	if *printFlag {
		// Сценарий ./myapp -o
		printScript(*limitFlag)
	} else {
		// Сценарий ./myapp
		downloadScript()
	}
}

func downloadScript() {
	err := requests.DBDownloadComics(app, 1, base_batch_size)
	if err != nil {
		log.Fatal(err)
	}
}

func printScript(limit int) {
	err := requests.DBPrintComics(app, limit)
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
