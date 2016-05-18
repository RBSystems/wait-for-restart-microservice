#!/usr/bin/env bash
# Pasted into Jenkins to build (will eventually be fleshed out to work with a Docker Hub and Amazon AWS)

echo "Stopping running application"
docker stop wait-for-reboot-microservice
docker rm wait-for-reboot-microservice

echo "Building container"
docker build -t byuoitav/wait-for-reboot-microservice .

echo "Starting the new version"
docker run -d --restart=always --name wait-for-reboot-microservice -p 8003:8003 byuoitav/wait-for-reboot-microservice:latest
