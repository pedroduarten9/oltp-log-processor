# OLTP log processor

This service will be responsible for printing logs for a given key on a given duration.

## Development

The development docs are [here](./docs/development.md).
It encompasses documentation that helps extend this service and facilitates its operation for developers.

## Decisions

The decisions taken along the challenge will be documented [here](./docs/decisions.md)

## How to run (via Golang)

To run the service you can either use Golang directly or use the helper Makefile, the commands for each of the approaches are below.

### Commands

`go run main.go serve` or
`make serve` 

### Prerequisites

Have Golang installed (version >= 1.23)

## How to run (via Docker)

There is a Dockerfile in case you don't have Golang installed. The command has the same flags and usages as [executing via Golang](#how-to-run-via-golang)
First you need to build the image.   

`docker build -t <tag> .`  

### Command

`docker run <tag> serve` 

### Prerequisites

Have Docker installed

## How to run (via shell scripts)
  
## How to start the processor

In order to start the server a helper file was created `run.sh`.runs the container with the application, when finished it removes the container.  

## How to check the logs

In order to check the logs of the application one can run `inspect.sh`. Once executed it will be continuously listen to the logs on the container.

## How to delete the resources created

In order to delete the resources created a helper file was created `cleanup.sh`. This executable deletes the image.

### Prerequisites

Have Docker installed