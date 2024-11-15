#!/bin/bash

# Default port if no argument is provided
PORT=${1:-8000}

# Run the Docker container with the dynamic port
sudo docker run -p $PORT:8000 --rm -v $(pwd):/app -v /app/tmp --name kalasya-air kalaysa
