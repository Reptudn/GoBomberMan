package routes

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func JoinGame(c *gin.Context, kubeClient kubernetes.Interface) {

	gameId := c.Param("gameId")
	if gameId == "" {
		c.JSON(404, gin.H{"status": "No game id supplied!"})
		return
	}

	serviceName := fmt.Sprintf("game-service-%s", gameId)
	service, err := kubeClient.CoreV1().Services("default").Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		c.JSON(404, gin.H{"status": "error", "message": "Game not found!"})
		return
	}

	if len(service.Spec.Ports) == 0 {
		fmt.Printf("Service %s has no ports defined\n", serviceName)
		c.JSON(404, gin.H{"status": "error", "message": "Game service has no ports"})
		return
	}

	nodePort := service.Spec.Ports[0].NodePort
	fmt.Printf("Found game service. NodePort: %d\n", nodePort)

	c.JSON(200, gin.H{
		"status":      "success",
		"gameId":      gameId,
		"url":         fmt.Sprintf("%s:%d", "localhost", nodePort),
		"servicePort": nodePort,
	})
}
