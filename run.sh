#!/bin/bash

# Сборка основного сервиса
echo "Building main service"
docker build -t task-tracker-back:v1 .
cd ..

# # Сборка почтового сервиса
echo "Building email service"
cd task-tracker-email-sender
docker build -t task-tracker-email-sender:v1 .
cd ..

# # Сборка сервиса планировщика 
echo "Building scheduler service"
cd task-tracker-scheduler
docker build -t task-tracker-scheduler:v1 .
cd ..

# Перезапуск docker compose
echo "Restarting docker compose..."
cd task-tracker-compose
docker compose down
docker compose up -d

echo "Docker run complete!"