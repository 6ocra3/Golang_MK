package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"makar/stemmer/pkg/config"
	"makar/stemmer/pkg/words"
	"net/http"
	"time"
)

type RawComic struct {
	Url        string `json:"img"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
}

type Comic struct {
	Url      string
	Keywords []string
}

var client = http.Client{Timeout: 30 * time.Second}
var configData *config.Config

func Init(config *config.Config) error {
	configData = config
	return nil
}

func DownloadComics(start int, end int) (map[int]*Comic, error) {
	fmt.Printf("Start fetching %d %d\n", start, end)
	comics := make(map[int]*Comic)
	fmt.Print("Fetched: ")

	comicsChan := make(chan *Comic, end-start+1)
	errChan := make(chan error, end-start+1)

	// Считывание комиксов по айти с помощью горутин
	for i := start; i <= end; i++ {
		comic, err := DownloadComic(i)
		if err != nil {
			return nil, err
		}
		comics[i] = comic
	}

	// Записывание каждого комикса в map
	for i := start; i <= end; i++ {
		select {
		case comic := <-comicsChan:
			comics[i] = comic
		case err := <-errChan:
			return nil, err
		}
	}

	fmt.Print("\n")
	fmt.Printf("Finish fetching %d %d\n", start, end)
	return comics, nil

}

func DownloadComic(id int) (*Comic, error) {

	// Считывание одного комикса

	url := fmt.Sprintf("%s/%d/info.0.json", configData.SourceURL, id)
	resp, err := client.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var comic RawComic

	json.Unmarshal(body, &comic)

	// Преобразование комикса в нужный формат. Стеминг, фильтр полей

	description := comic.Transcript + comic.Alt

	stemmedDescription, err := words.StemmString(description)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	processedComic := Comic{comic.Url, stemmedDescription}

	return &processedComic, nil
}
