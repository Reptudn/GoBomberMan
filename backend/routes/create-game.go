package routes

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func CreateGame(c *gin.Context, kubeClient kubernetes.Interface) {

	gameId := uuid.New().String()

	if err := createGamePod(gameId, kubeClient); err != nil {
		c.JSON(500, gin.H{"status": "error creating game pod"})
		return
	}

	servicePort, err := createGameService(gameId, kubeClient)
	if err != nil {
		c.JSON(500, gin.H{"status": "error creating game service"})
		return
	}

	fmt.Printf("Game service created for game %s on port %d\n", gameId, servicePort)

	c.JSON(200, gin.H{
		"status":      "Game created successfully",
		"gameId":      gameId,
		"socketUrl":   fmt.Sprintf("ws://%s:%d/ws", "localhost", servicePort),
		"servicePort": servicePort,
	})

}

func createGameService(gameId string, kubeClient kubernetes.Interface) (int32, error) {
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("game-service-%s", gameId),
			Namespace: "default",
			Labels: map[string]string{
				"app":    "bomberman-game-server",
				"gameId": gameId,
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Selector: map[string]string{
				"gameId": gameId,
			},
			Ports: []corev1.ServicePort{
				{
					Port:       8080,
					TargetPort: intstr.FromInt(8080),
					Name:       "websocket",
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}

	createdService, err := kubeClient.CoreV1().Services("default").Create(context.TODO(), &service, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error creating service for game %s: %v\n", gameId, err)
		return 0, err
	}

	nodePort := createdService.Spec.Ports[0].NodePort
	fmt.Printf("Service created for game %s on node port %d\n", gameId, nodePort)

	return nodePort, nil
}

func createGamePod(gameId string, kubeClient kubernetes.Interface) error {
	pod := corev1.Pod{ObjectMeta: metav1.ObjectMeta{
		Name:      fmt.Sprintf("game-server-%s", gameId),
		Namespace: "default",
		Labels: map[string]string{
			"app":    "bomberman-game-server",
			"gameId": gameId,
		},
	},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "game-container",
					Image:           "bomberman-game-image:latest",
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

	_, err := kubeClient.CoreV1().Pods("default").Create(context.TODO(), &pod, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error creating pod for game %s: %v\n", gameId, err)
		return err
	} else {
		fmt.Printf("Pod created for game %s\n", gameId)
	}
	return nil
}
