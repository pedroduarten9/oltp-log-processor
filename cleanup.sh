#!/bin/bash
docker stop log-processor
docker image rm log-processor
echo "Docker image 'log-processor' removed successfully."