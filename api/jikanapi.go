// webserver.go
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

type Image struct {
	ImageURL      string `json:"image_url"`
	SmallImageURL string `json:"small_image_url"`
	LargeImageURL string `json:"large_image_url"`
}

type Anime struct {
	MALID    int    `json:"mal_id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Synopsis string `json:"synopsis"`
	Images   struct {
		JPG Image `json:"jpg"`
	} `json:"images"`
}

type JikanResponse struct {
	Data []Anime `json:"data"`
}

func searchAnimeJikan(searchQuery string) (*JikanResponse, error) {
	encodedQuery := url.QueryEscape(searchQuery)
	apiURL := fmt.Sprintf("https://api.jikan.moe/v4/anime?q=%s", encodedQuery)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", resp.Status)
	}

	var results JikanResponse
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		query := r.FormValue("query")
		if query == "" {
			http.Error(w, "Query parameter is missing", http.StatusBadRequest)
			return
		}

		results, err := searchAnimeJikan(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles("w/index.html"))
		err = tmpl.Execute(w, results)
		if err != nil {
			http.Error(w, "Failed to render results", http.StatusInternalServerError)
			return
		}
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
