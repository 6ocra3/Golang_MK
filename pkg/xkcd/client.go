package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RawComic struct {
	ID         int
	Url        string `json:"img"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
}

type Client struct {
	SourceURL string
}

var httpClient = http.Client{Timeout: 30 * time.Second}

func Init(SourceURL string) (*Client, error) {
	client := &Client{
		SourceURL: SourceURL,
	}
	return client, nil
}

func DownloadComics(client *Client, start int, end int) ([]*RawComic, error) {
	fmt.Printf("Start fetching %d %d\n", start, end)
	comics := make([]*RawComic, 0, end-start)
	fmt.Print("Fetched: ")

	// Считывание комиксов по айди
	for i := start; i <= end; i++ {
		comic, err := DownloadComic(client, i)
		if err != nil {
			return nil, err
		}
		fmt.Printf("%d ", i)
		comics = append(comics, comic)
		// comics[i] = comic
	}

	fmt.Print("\n")
	fmt.Printf("Finish fetching %d %d\n", start, end)
	return comics, nil

}

func DownloadComic(client *Client, id int) (*RawComic, error) {

	// Считывание одного комикса
	url := fmt.Sprintf("%s/%d/info.0.json", client.SourceURL, id)
	resp, err := httpClient.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var comic RawComic

	json.Unmarshal(body, &comic)
	comic.ID = id

	return &comic, nil
}
