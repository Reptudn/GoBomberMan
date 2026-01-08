package routes

import (
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

func JoinGame(c *gin.Context, kubeClient kubernetes.Interface) {

	gameId, exists := c.Params.Get("gameId")
	if !exists {
		c.JSON(500, gin.H{"status": "No game id supplied!"})
	}

	kubeClient.CoreV1().Pods("default").List()

	c.JSON(200, gin.H{"status": "join game endpoint"})
}
