package routes

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateGame(c *gin.Context, kubeClient kubernetes.Interface) {

	gameId := uuid.New().String()


	if createGamePod(gameId, kubeClient) != nil {
		c.JSON(500, gin.H{"status": "error creating game pod"})
		return
	}

	c.JSON(200, gin.H{"status": "game created", "gameId": gameId})

}

func createGamePod(gameId string, kubeClient kubernetes.Interface) error {
	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{
		Name:      fmt.Sprintf("game-server-%s", gameId),
		Namespace: "default",
		Labels: map[string]string{
			"app": "bomberman-game-server",
			"gameId": gameId,
		},
	},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "game-container",
					Image: "bomberman-game-image:latest",
					ImagePullPolicy: corev1.PullNever,
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 8080,
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	_, err := kubeClient.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error creating pod for game %s: %v\n", gameId, err)
		return err
	} else {
		fmt.Printf("Pod created for game %s\n", gameId)
	}
	return nil
}