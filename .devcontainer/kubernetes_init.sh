#/bin/sh

# Fix docker permissions
sudo chmod 666 /var/run/docker.sock

# Create kubernetes directory if it doesn't exist
mkdir -p /workspaces/StackItManGO/kubernetes

# Create minikube cluster if it doesn't exist with reduced memory requirements
minikube status > /dev/null 2>&1 || minikube start --driver=docker --memory=1500mb --cpus=2

# Update kubectl context to use minikube
minikube update-context

# Export minikube's kubeconfig to the kubernetes directory
cp ~/.kube/config ../kubernetes/config.yaml
chmod 600 ../kubernetes/config.yaml

# Export KUBECONFIG so it's accessible
export KUBECONFIG=$(pwd)/../kubernetes/config.yaml

echo "Kubernetes cluster created and config saved to kubernetes/config.yaml"