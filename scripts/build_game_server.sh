#/bin/sh

cd ../game-server

docker build -t game-server:latest .

minikube image load game-server:latest