package routes

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func JoinGame(c echo.Context, kubeClient kubernetes.Interface) error {

	gameId := c.Param("gameId")
	if gameId == "" {
		return c.JSON(404, map[string]interface{}{"status": "No game id supplied!"})
	}

	serviceName := fmt.Sprintf("game-service-%s", gameId)
	service, err := kubeClient.CoreV1().Services("default").Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return c.JSON(404, map[string]interface{}{"status": "error", "message": "Game not found!"})
	}

	if len(service.Spec.Ports) == 0 {
		fmt.Printf("Service %s has no ports defined\n", serviceName)
		return c.JSON(404, map[string]interface{}{"status": "error", "message": "Game service has no ports"})
	}

	nodePort := service.Spec.Ports[0].NodePort
	fmt.Printf("Found game service. NodePort: %d\n", nodePort)

	return c.JSON(200, map[string]interface{}{
		"status":      "success",
		"gameId":      gameId,
		"url":         fmt.Sprintf("%s:%d", "localhost", nodePort),
		"servicePort": nodePort,
	})
}
