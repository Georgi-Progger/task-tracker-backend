#!/bin/bash

# Сборка основного сервиса
echo "Building main service..."
docker build -t task-tracker-back:v1 .
cd ..

# Перезапуск docker compose
echo "Restarting docker compose..."
cd task-tracker-compose
docker compose down
docker compose up -d

echo "Docker run complete!"