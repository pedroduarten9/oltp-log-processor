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