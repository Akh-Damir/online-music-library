package handlers

import (
	"net/http"
	"online-music-library/models"
	"online-music-library/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Get Songs
// @Description Get a list of songs with optional filters
// @Tags Songs
// @Param group query string false "Filter by group name"
// @Param song query string false "Filter by song name"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {array} models.Song
// @Failure 500 {object} map[string]string
// @Router /songs [get]
func GetSongs(c *gin.Context) {
	group := c.Query("group") // Фильтрация по группе
	song := c.Query("song")   // Фильтрация по названию песни
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Получение данных из модели
	songs, err := models.FetchSongs(group, song, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch songs"})
		return
	}

	c.JSON(http.StatusOK, songs)
}

// @Summary Get Song Text
// @Description Retrieve the text of a song by its ID, with pagination by verses
// @Tags Songs
// @Param id path int true "Song ID"
// @Param page query int false "Page number" default(1)
// @Success 200 {array} string "Verses of the song"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /songs/{id}/text [get]
func GetSongText(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 1 // По одному куплету за раз

	text, err := models.FetchSongText(id, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song text"})
		return
	}

	c.JSON(http.StatusOK, text)
}

// @Summary Add Song
// @Description Add a new song to the library
// @Tags Songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Song details"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /songs [post]
func AddSong(c *gin.Context) {
	var song models.Song

	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	details, err := utils.FetchSongDetails(song.GroupName, song.SongName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song details"})
		return
	}

	modelDetails := models.ConvertUtilsToModelsDetails(details)

	if err := models.SaveSong(&song, modelDetails); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save song"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Song added successfully"})
}

// @Summary Update Song
// @Description Update the details of an existing song by its ID
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Updated song details"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /songs/{id} [put]
func UpdateSong(c *gin.Context) {
	id := c.Param("id")

	var updatedSong models.Song
	if err := c.ShouldBindJSON(&updatedSong); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := models.UpdateSong(id, &updatedSong); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}

// @Summary Delete Song
// @Description Delete a song by its ID
// @Tags Songs
// @Param id path int true "Song ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /songs/{id} [delete]

func DeleteSong(c *gin.Context) {
	id := c.Param("id")

	if err := models.DeleteSong(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}
