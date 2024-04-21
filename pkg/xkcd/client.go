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

var httpClient = http.Client{
	Timeout: 5 * time.Second,
}

func Init(SourceURL string) (*Client, error) {
	client := &Client{
		SourceURL: SourceURL,
	}
	return client, nil
}

func DownloadComic(client *Client, id int) (*RawComic, error) {
	// Считывание одного комикса
	url := fmt.Sprintf("%s/%d/info.0.json", client.SourceURL, id)
	var err error
	var resp *http.Response
	for i := 0; i < 5; i++ {
		resp, err = httpClient.Get(url)

		if err == nil && resp.StatusCode == http.StatusOK || id == 404 {
			break
		}

		//fmt.Printf("%d-error ", id)

	}

	if err != nil {
		return nil, err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK && id != 404 {
		return nil, fmt.Errorf("failed to fetch comic %d: %s", id, resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)

	var comic RawComic

	json.Unmarshal(body, &comic)
	comic.ID = id

	return &comic, nil
}
