#!/bin/bash
set -e

docker run --rm -d --name log-processor log-processor serve "$@"