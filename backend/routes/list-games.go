package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type GameInfo struct {
	Name           string `json:"name"`
	GameId         string `json:"gameId"`
	CurrentPlayers int    `json:"currentPlayers"`
	MaxPlayers     int    `json:"maxPlayers"`
}

type gameListResponse struct {
	Status string     `json:"status"`
	Count  int        `json:"count"`
	Games  []GameInfo `json:"games"`
}

func ListGames(c *gin.Context, kubeClient kubernetes.Interface) {
	pods, err := kubeClient.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{LabelSelector: "app=bomberman-game-server"})

	if err != nil {
		c.JSON(500, gin.H{"status": "error listing game pods"})
		return
	}

	var games []GameInfo
	for _, pod := range pods.Items {
		var gameInfo GameInfo

		gameInfo.CurrentPlayers = 0 // Placeholder for actual player count logic
		gameInfo.MaxPlayers = 4     // Assuming max players is 4 for all games
		gameInfo.Name = pod.Name
		gameInfo.GameId = pod.Labels["gameId"]

		// TODO: Implement logic to fetch current player count from pod status
		// TODO: Only add the game to the array when the container is ready
		if pod.Status.Phase == "Running" {
			games = append(games, gameInfo)
		}
	}

	c.JSON(200, gin.H{
		"status": "success",
		"count":  len(games),
		"games":  games,
	})

}
