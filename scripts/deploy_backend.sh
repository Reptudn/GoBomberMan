#!/bin/bash

# Deploy Backend (Manager Node) to Kubernetes
# This script builds the Docker image and deploys it to Docker Desktop Kubernetes

set -e  # Exit on any error

# Get the absolute path to the project root (parent of scripts directory)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "=========================================="
echo "Deploying Backend to Kubernetes"
echo "=========================================="
echo "Project root: $PROJECT_ROOT"

# Step 1: Check if kubectl is available
echo ""
echo "[1/5] Checking if kubectl is installed..."
if ! command -v kubectl &> /dev/null; then
    echo "Error: kubectl is not installed or not in PATH"
    exit 1
fi
echo "✓ kubectl found"

# Step 2: Check if Kubernetes cluster is running
echo ""
echo "[2/5] Checking if Kubernetes cluster is accessible..."
if ! kubectl cluster-info &> /dev/null; then
    echo "Error: Cannot connect to Kubernetes cluster"
    echo "Make sure Docker Desktop is running and Kubernetes is enabled"
    exit 1
fi
echo "✓ Kubernetes cluster is accessible"

# Step 3: Build the Docker image
echo ""
echo "[3/5] Building backend Docker image..."
cd "$PROJECT_ROOT/backend"
docker build -t go-bomberman-backend:latest .
if [ $? -ne 0 ]; then
    echo "Error: Failed to build Docker image"
    exit 1
fi
echo "✓ Docker image built successfully"

# Step 4: Apply RBAC configuration
echo ""
echo "[4/5] Applying RBAC configuration..."
cd "$PROJECT_ROOT/kubernetes"
kubectl apply -f backend-rbac.yaml
if [ $? -ne 0 ]; then
    echo "Error: Failed to apply RBAC configuration"
    exit 1
fi
echo "✓ RBAC configuration applied"

# Step 5: Deploy the backend
echo ""
echo "[5/5] Deploying backend to Kubernetes..."
kubectl apply -f backend-deployment.yaml
if [ $? -ne 0 ]; then
    echo "Error: Failed to deploy backend"
    exit 1
fi
echo "✓ Backend deployed successfully"

# Wait for deployment to be ready
echo ""
echo "Waiting for backend to be ready..."
kubectl wait --for=condition=available --timeout=60s deployment/backend

# Show deployment status
echo ""
echo "=========================================="
echo "Deployment Complete!"
echo "=========================================="
echo ""
echo "Backend Status:"
kubectl get pods -l app=backend
echo ""
echo "Service Status:"
kubectl get service backend-service
echo ""
echo "Backend should be accessible at: http://localhost:8080"
echo ""
echo "Useful commands:"
echo "  - View logs: kubectl logs -l app=backend -f"
echo "  - Check status: kubectl get pods"
echo "  - Restart: kubectl rollout restart deployment/backend"
echo "  - Delete: kubectl delete -f $PROJECT_ROOT/kubernetes/backend-deployment.yaml"
echo ""
