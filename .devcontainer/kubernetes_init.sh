#!/bin/sh

# Fix docker permissions if needed
sudo chmod 666 /var/run/docker.sock

WORKSPACE_DIR="$(cd "$(dirname "$0")/.." && pwd)"
KUBE_DIR="$WORKSPACE_DIR/kubernetes"

# Ensure kubernetes directory exists in the current workspace
mkdir -p "$KUBE_DIR"

# Create minikube cluster if it doesn't exist with sufficient memory
if ! minikube status >/dev/null 2>&1; then
	minikube start --driver=docker --memory=2000mb --cpus=2
fi

# Update kubectl context to use minikube
minikube update-context

# Export minikube's kubeconfig to the kubernetes directory
if [ -f "$HOME/.kube/config" ]; then
	cp "$HOME/.kube/config" "$KUBE_DIR/config.yaml"
	chmod 600 "$KUBE_DIR/config.yaml"
	export KUBECONFIG="$KUBE_DIR/config.yaml"
	echo "Kubernetes cluster created and config saved to kubernetes/config.yaml"
else
	echo "minikube kubeconfig not found at $HOME/.kube/config; did minikube start succeed?" >&2
	exit 1
fi