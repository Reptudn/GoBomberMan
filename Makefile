# GoBomberMan Makefile
# Convenience commands for building and deploying the project

.PHONY: help build-all build-backend build-frontend build-game-server deploy-backend k8s-status k8s-logs k8s-restart k8s-clean test-backend clean-games

# Get absolute path to project root
ROOT_DIR := $(shell pwd)

# Default target - show help
help:
	@echo "GoBomberMan - Available Commands"
	@echo "================================="
	@echo ""
	@echo "Building:"
	@echo "  make build-all          - Build all components (frontend, backend, game-server)"
	@echo "  make build-backend      - Build only the backend Docker image"
	@echo "  make build-frontend     - Build only the frontend Docker image"
	@echo "  make build-game-server  - Build only the game-server Docker image"
	@echo ""
	@echo "Deployment:"
	@echo "  make deploy-backend     - Deploy backend to Kubernetes (builds + deploys)"
	@echo "  make deploy-quick       - Quick deploy (skips build, just applies k8s configs)"
	@echo ""
	@echo "Kubernetes Management:"
	@echo "  make k8s-status         - Show status of all pods and services"
	@echo "  make k8s-logs           - Show backend logs (follows in real-time)"
	@echo "  make k8s-restart        - Restart the backend deployment"
	@echo "  make k8s-clean          - Remove all deployed resources"
	@echo "  make k8s-describe       - Detailed pod information"
	@echo "  make clean-games        - Delete all game server pods and services"
	@echo ""
	@echo "Testing:"
	@echo "  make test-backend       - Test backend endpoints"
	@echo "  make test-ping          - Quick ping test"
	@echo "  make create-game        - Create a test game"
	@echo ""
	@echo "Development:"
	@echo "  make rebuild-backend    - Rebuild backend and restart in k8s"
	@echo "  make rebuild-game-server - Rebuild game server image"
	@echo ""

# Build all components
build-all: build-backend build-frontend build-game-server
	@echo "✓ All components built successfully!"

# Build individual components
build-backend:
	@echo "Building backend..."
	@cd $(ROOT_DIR)/backend && docker build -t go-bomberman-backend:latest .
	@echo "✓ Backend built successfully"

build-frontend:
	@echo "Building frontend..."
	@cd $(ROOT_DIR)/frontend && docker build -t go-bomberman-frontend:latest .
	@echo "✓ Frontend built successfully"

build-game-server:
	@echo "Building game server..."
	@cd $(ROOT_DIR)/game-server && docker build -t go-bomberman-game-server:latest .
	@echo "✓ Game server built successfully"

# Deploy backend to Kubernetes (full deployment)
deploy-backend: build-backend
	@echo "Applying RBAC configuration..."
	@kubectl apply -f $(ROOT_DIR)/kubernetes/backend-rbac.yaml
	@echo "Deploying backend..."
	@kubectl apply -f $(ROOT_DIR)/kubernetes/backend-deployment.yaml
	@echo "Waiting for deployment to be ready..."
	@kubectl wait --for=condition=available --timeout=60s deployment/backend || true
	@echo "✓ Deployment complete!"

# Quick deploy (assumes images are already built)
deploy-quick:
	@echo "Applying Kubernetes configurations..."
	@kubectl apply -f $(ROOT_DIR)/kubernetes/backend-rbac.yaml
	@kubectl apply -f $(ROOT_DIR)/kubernetes/backend-deployment.yaml
	@echo "Waiting for deployment to be ready..."
	@kubectl wait --for=condition=available --timeout=60s deployment/backend || true
	@echo "✓ Deployment complete!"

# Show Kubernetes status
k8s-status:
	@echo "=== Pods ==="
	@kubectl get pods
	@echo ""
	@echo "=== Services ==="
	@kubectl get services
	@echo ""
	@echo "=== Deployments ==="
	@kubectl get deployments
	@echo ""
	@echo "=== Game Servers ==="
	@kubectl get pods -l app=bomberman-game-server

# View backend logs
k8s-logs:
	@echo "Following backend logs (Ctrl+C to exit)..."
	@kubectl logs -l app=backend -f

# Show last 50 log lines
k8s-logs-tail:
	@kubectl logs -l app=backend --tail=50

# View game server logs
k8s-logs-games:
	@echo "Game server logs:"
	@kubectl logs -l app=bomberman-game-server --tail=50

# Restart backend deployment
k8s-restart:
	@echo "Restarting backend deployment..."
	@kubectl rollout restart deployment/backend
	@kubectl rollout status deployment/backend

# Clean up all Kubernetes resources
k8s-clean:
	@echo "Removing all deployed resources..."
	-@kubectl delete -f $(ROOT_DIR)/kubernetes/backend-deployment.yaml
	-@kubectl delete -f $(ROOT_DIR)/kubernetes/backend-rbac.yaml
	@echo "✓ Cleanup complete!"

# Clean up game servers
clean-games:
	@echo "Deleting all game server pods and services..."
	-@kubectl delete pods -l app=bomberman-game-server
	-@kubectl delete services -l app=bomberman-game-server
	@echo "✓ Game servers cleaned up!"

# Describe pods for debugging
k8s-describe:
	@echo "=== Backend Pod Details ==="
	@kubectl describe pod -l app=backend

# Test backend endpoints
test-backend: test-ping test-list-games

test-ping:
	@echo "Testing /ping endpoint..."
	@curl -s http://localhost:8080/ping | jq . 2>/dev/null || curl -s http://localhost:8080/ping || echo "Backend not responding"

test-list-games:
	@echo "Testing /list-games endpoint..."
	@curl -s http://localhost:8080/list-games | jq . 2>/dev/null || curl -s http://localhost:8080/list-games || echo "Backend not responding"

create-game:
	@echo "Creating a new game..."
	@curl -s -X POST http://localhost:8080/create-game | jq . 2>/dev/null || curl -s -X POST http://localhost:8080/create-game

# Rebuild backend and restart in Kubernetes
rebuild-backend: build-backend k8s-restart
	@echo "✓ Backend rebuilt and restarted!"

# Rebuild game server
rebuild-game-server: build-game-server clean-games
	@echo "✓ Game server rebuilt! Old game pods deleted."

# Check if Kubernetes is running
check-k8s:
	@kubectl cluster-info > /dev/null 2>&1 || (echo "Error: Kubernetes cluster not accessible. Is Docker Desktop running?" && exit 1)
	@echo "✓ Kubernetes cluster is accessible"

# Port forward backend (useful if LoadBalancer isn't working)
port-forward:
	@echo "Port forwarding backend to localhost:8080 (Ctrl+C to stop)..."
	@kubectl port-forward deployment/backend 8080:8080

# Show all resources with backend label
k8s-all:
	@kubectl get all -l app=backend

# Show all game resources
k8s-games:
	@kubectl get pods,services -l app=bomberman-game-server

# Watch pod status (useful during deployment)
watch-pods:
	@watch -n 1 kubectl get pods

# Full reset and redeploy
reset: k8s-clean clean-games deploy-backend
	@echo "✓ Full reset complete!"
