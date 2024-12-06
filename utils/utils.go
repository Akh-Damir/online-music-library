package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SongDetails struct {
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func FetchSongDetails(group, song string) (*SongDetails, error) {
	apiURL := fmt.Sprintf("%s?group=%s&song=%s", os.Getenv("EXTERNAL_API_URL"), group, song)
	fmt.Printf("Requesting URL: %s\n", apiURL)
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return nil, fmt.Errorf("failed to fetch details: %w", err)
	}
	defer resp.Body.Close()

	var details SongDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &details, nil
}
