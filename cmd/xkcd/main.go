package main

import (
	"flag"
	"log"
	"makar/stemmer/pkg/config"
	"makar/stemmer/pkg/database"
	"makar/stemmer/pkg/requests"
	"makar/stemmer/pkg/xkcd"
)

const base_batch_size = 50

func main() {

	printFlag := flag.Bool("o", false, "Флаг для вывода скаченных комиксов")
	limitFlag := flag.Int("n", -1, "Флаг для установки лимита показываемых комиксов")

	flag.Parse()

	const configPath = "config.yaml"
	config, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	initAll(config)
	if *printFlag {
		printScript(-1)
	} else if *limitFlag != -1 {
		printScript(*limitFlag)
	} else {
		downloadScript()
	}
}

func downloadScript() {
	err := requests.DBDownloadComics(1, base_batch_size)
	if err != nil {
		log.Fatal(err)
	}
}

func printScript(limit int) {
	err := requests.DBPrintComics(limit)
	if err != nil {
		log.Fatal(err)
	}

}

func initAll(config *config.Config) {
	err := database.Init(config)
	if err != nil {
		log.Fatal(err)
	}

	err = xkcd.Init(config)
	if err != nil {
		log.Fatal(err)
	}
}
