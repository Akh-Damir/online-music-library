package models

import (
	"fmt"
	"online-music-library/utils"
	"strings"
)

type Song struct {
	ID          int    `json:"id"`
	GroupName   string `json:"group_name"`
	SongName    string `json:"song_name"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
type SongDetails struct {
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func ConvertUtilsToModelsDetails(details *utils.SongDetails) *SongDetails {
	return &SongDetails{
		ReleaseDate: details.ReleaseDate,
		Text:        details.Text,
		Link:        details.Link,
	}
}

func FetchSongs(group, song string, page, limit int) ([]Song, error) {
	var songs []Song
	offset := (page - 1) * limit

	query := `
		SELECT id, group_name, song_name, release_date, text, link 
		FROM songs WHERE 1=1
	`
	if group != "" {
		query += fmt.Sprintf(" AND group_name ILIKE '%%%s%%'", group)
	}
	if song != "" {
		query += fmt.Sprintf(" AND song_name ILIKE '%%%s%%'", song)
	}
	query += " ORDER BY id LIMIT $1 OFFSET $2"

	rows, err := DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s Song
		if err := rows.Scan(&s.ID, &s.GroupName, &s.SongName, &s.ReleaseDate, &s.Text, &s.Link); err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}

	return songs, nil
}

func FetchSongText(songID string, page, limit int) ([]string, error) {
	var fullText string
	err := DB.QueryRow("SELECT text FROM songs WHERE id = $1", songID).Scan(&fullText)
	if err != nil {
		return nil, err
	}

	verses := splitVerses(fullText)
	start := (page - 1) * limit
	end := start + limit
	if start >= len(verses) {
		return nil, nil
	}
	if end > len(verses) {
		end = len(verses)
	}

	return verses[start:end], nil
}

func splitVerses(text string) []string {
	return strings.Split(text, "\n\n")
}

func SaveSong(song *Song, details *SongDetails) error {
	query := `
		INSERT INTO songs (group_name, song_name, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := DB.Exec(query, song.GroupName, song.SongName, details.ReleaseDate, details.Text, details.Link)
	return err
}
func UpdateSong(id string, song *Song) error {
	query := `
		UPDATE songs 
		SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5
		WHERE id = $6
	`
	_, err := DB.Exec(query, song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link, id)
	return err
}
func DeleteSong(id string) error {
	query := "DELETE FROM songs WHERE id = $1"
	_, err := DB.Exec(query, id)
	return err
}
