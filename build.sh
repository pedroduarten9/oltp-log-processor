#!/bin/bash
set -e

docker build -t log-processor .
echo "Docker image 'log-processor' built successfully."