package main

import (
	"bomberman-backend/routes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	fmt.Println("Starting Backend Server...")

	// Try KUBECONFIG env var first, then fall back to workspace config
	kubeconfig, exists := os.LookupEnv("KUBECONFIG")
	if !exists || kubeconfig == "" {
		kubeconfig = filepath.Join("..", "kubernetes", "config.yaml")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/list-games", func(c *gin.Context) {
		routes.ListGames(c, kubeClient)
	})

	r.POST("/create-game", func(c *gin.Context) {
		routes.CreateGame(c, kubeClient)
	})

	r.POST("/join-game", func(c *gin.Context) {
		routes.JoinGame(c, kubeClient)
	})

	r.Run()
}
