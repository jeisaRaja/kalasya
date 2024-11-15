#!/bin/bash

# Stop the Docker Compose services
echo "Stopping Docker Compose services..."

# Use docker-compose to stop the services
sudo docker compose down

# Optional: Remove all stopped containers, unused networks, and volumes
# docker-compose down --volumes --remove-orphans

echo "Docker Compose services stopped."
