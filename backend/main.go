package main

import (
	"bomberman-backend/routes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	fmt.Println("Starting Backend Server...")

	var config *rest.Config
	var err error

	// Try in-cluster config first (when running in Kubernetes)
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig file (for local development)
		fmt.Println("Not running in cluster, using kubeconfig file...")
		kubeconfig, exists := os.LookupEnv("KUBECONFIG")
		if !exists || kubeconfig == "" {
			kubeconfig = filepath.Join("..", "kubernetes", "config.yaml")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	} else {
		fmt.Println("Running in cluster, using in-cluster config")
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	r := gin.Default()

	r.Use(CORSMiddleware())

	r.GET("/status", func(c *gin.Context) {
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

	r.POST("/join-game/:gameId", func(c *gin.Context) {
		routes.JoinGame(c, kubeClient)
	})

	fmt.Println("Backend Server started")
	r.Run()
}
