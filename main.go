package main

import (
	_ "embed"
	"flag"
	"fmt"
	"moodtheme/api/middleware"
	"moodtheme/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

type song struct {
	Song string `json:"song" binding:"required"`
}

func main() {

	filepath := flag.String("themes", "./example_themes.json", "Absolute path to themes.json")

	if err := data.LoadData(*filepath); err != nil {
		fmt.Println("Warning: Cannot read json themes", err)
	}

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	router.POST("/theme", middleware.AntiSpamMiddleware(), func(c *gin.Context) {
		var validatedSong song
		if err := c.ShouldBindJSON(&validatedSong); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "The song name should't be empty"})
			return
		}
		if err := data.BroadcastTheme(validatedSong.Song); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Something went wrong: %s", err)})
		} else {
			value, _ := data.FetchTheme(validatedSong.Song)
			c.JSON(http.StatusAccepted, gin.H{"message": fmt.Sprintf("Changed accepted! Theme modified to: %s", value)})
		}
	})

	router.Run(":8080")

}
