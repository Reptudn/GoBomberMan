package main

import (
	"bomberman-backend/routes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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

	e := echo.New()

	isProd := os.Getenv("ENV") == "production"

	if isProd {
		e.Debug = false
		e.HideBanner = true
		e.HidePort = true
	} else {
		e.Debug = true
		e.Use(middleware.RequestLogger())
	}

	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"Content-Length", "Accept-Encoding", "X-CSRF-Token", "Cache-Control", "X-Requested-With",
		},
		AllowCredentials: true,
	}))

	// TODO: change responses to JSONP https://echo.labstack.com/docs/cookbook/jsonp#server
	e.GET("/ping", routes.Ping)
	e.RouteNotFound("/*", routes.RouteNotFound)

	e.GET("/list-games", func(c echo.Context) error {
		return routes.ListGames(c, kubeClient)
	})

	e.POST("/create-game", func(c echo.Context) error {
		return routes.CreateGame(c, kubeClient)
	})

	e.POST("/join-game/:gameId", func(c echo.Context) error {
		return routes.JoinGame(c, kubeClient)
	})

	fmt.Println("Backend Server started on port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
