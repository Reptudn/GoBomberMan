#/bin/sh

# build image
docker build -t go-bomberman-report:latest .

# apply to cluster
kubectl apply -f ../kubernetes/report-deployment.yaml