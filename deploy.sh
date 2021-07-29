#!/bin/bash
cd "$(dirname "$0")"
#disable database local for using same port
sudo systemctl stop postgresql
#Deploy service
docker-compose up --force-recreate --build -d
# for removing old image
docker image prune -f
