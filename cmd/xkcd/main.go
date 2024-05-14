package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	adapters "makar/stemmer/adapters/http"
	"makar/stemmer/pkg/config"
	"makar/stemmer/pkg/database/json"
	database2 "makar/stemmer/pkg/database/sqlite"
	"makar/stemmer/pkg/requests"
	"makar/stemmer/pkg/xkcd"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	database2.InitSQLite("test")
}

func main2() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var port int
	flag.IntVar(&port, "p", 8080, "Флаг `-p` используется для ввода порта")
	flag.Parse()

	// Считывание конфига
	const configPath = "config.yaml"
	config, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Создание БД, подгрузка конфига в пакеты
	app := initAll(config)

	err = requests.DBDownloadComics(app, ctx, config.Parallel, config.IndexFile)
	if err != nil {
		fmt.Println(err)
	}

	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := requests.DBDownloadComics(app, ctx, config.Parallel, config.IndexFile)
				if err != nil {
					fmt.Println(err)
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	http.HandleFunc("/update", adapters.DBDownloadComicsAdapter(app, ctx, config.Parallel, config.IndexFile))
	http.HandleFunc("/pics", adapters.DBFindComicsAdapter(app, config.SearchLimit))

	serverAddress := ":" + strconv.Itoa(port)
	server := &http.Server{
		Addr: serverAddress,
	}

	go func() {
		fmt.Printf("Server is starting on port %d...\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("Shutting down gracefully, press Ctrl+C again to force")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exiting")

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
