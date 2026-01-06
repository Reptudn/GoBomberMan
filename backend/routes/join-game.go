package routes

import "github.com/gin-gonic/gin"

func JoinGame(c *gin.Context) {
	c.JSON(200, gin.H{"status": "join game endpoint"})
}