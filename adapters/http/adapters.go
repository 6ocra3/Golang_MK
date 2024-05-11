package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"makar/stemmer/pkg/requests"
	"net/http"
)

type DBDownloadComicsResponse struct {
	New   int `json:"new"`
	Total int `json:"total"`
}

type DBFindComicsResponse struct {
	Comics []string `json:"comics"`
}

func DBDownloadComicsAdapter(app *requests.App, ctx context.Context, parallel int, indexFile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		prevTotal := len(app.Db.Entries)

		err := requests.DBDownloadComics(app, ctx, parallel, indexFile)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to download comics: %v", err), http.StatusInternalServerError)
			return
		}

		newComics := len(app.Db.Entries) - prevTotal
		totalComics := len(app.Db.Entries)

		response := DBDownloadComicsResponse{
			newComics, totalComics,
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		}
	}
}

func DBFindComicsAdapter(app *requests.App, limit int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		request := r.URL.Query().Get("search")
		if request == "" {
			http.Error(w, "Missing search query", http.StatusBadRequest)
			return
		}

		err, comics := requests.DBFindComics(app, request, limit)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to find comics: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(comics); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		}
	}
}
