package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/songs", GetSongs)             // Получить список песен
	router.GET("/songs/:id/text", GetSongText) // Получить текст песни с пагинацией по куплетам
	router.POST("/songs", AddSong)             // Добавить новую песню
	router.PUT("/songs/:id", UpdateSong)       // Обновить песню
	router.DELETE("/songs/:id", DeleteSong)    // Удалить песню
}
