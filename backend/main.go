package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type game struct {
	ID   string
	Name string
}

var games = []game{};

func main() {
	fmt.Println("Starting Backend Server...")

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/list-games", func(c *gin.Context) {
		c.String(200, "Welcome to the Backend Server!")
	})

	r.POST("/create-game", func(c *gin.Context) {
		var newGame game
		if err := c.BindJSON(&newGame); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		games = append(games, newGame)
		c.JSON(200, gin.H{"status": "game created", "game": newGame})
	})

	r.POST("join-game", func(ctx *gin.Context) {
		
	})

	r.Run()
}