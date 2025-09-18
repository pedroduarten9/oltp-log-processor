# Decisions

This document will outline the technical and product decisions taken along the challenge.

## Tech Decisions

### Usage of Cobra

Cobra CLI is a CLI helper that makes helps with the running of the service, so I opted to use it.

### Usage of mocks

There are interfaces in case I wanted to mock requests, I saw no need for that since I can easily test the full flow without depending on external sources or flaky I/O. Therefore I decided not to use mocks.

### Support for log values

I opted to support multiple value types for the logs but I will only test for strings, otherwise it would take me much more time. I can also remove the support.

### Log processing response

For this iteration I can either find or not the key, I will not be processing errors, therefore the return is always the same.

### Concurrency

There are unused structs such as concurrent repo and concurrent service that use a synchronized map with sharding and go routines respectively, with the benchmarks result the decision was to keep the simpler versions.