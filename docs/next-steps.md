# Next steps

This document will document the PoC next steps, eventually transforming it from a PoC to a product.

## Input Validation

We don't have that much validations for the request apart from it being a gRPC defined struct, therefore we can introduce some business rules, with that we would be able to add the failed logs to the response.

Priority: P0

## Security

We should secure the server with TLS.

Priority: P0

## Monitoring

We should be able to monitor this service to understand it's health. We should check for it's load, validate it's data sources...

Priority: P0

## Load testing

Since this service should be able to handle a lot of request we should ensure that it does it's job, so we should load test it

Priority: P0

## Deployment on the cloud

This PoC will benefit a lot from being deployed on the cloud, we should consider using IaC for that.

Priority: P0

## Security

With the IaC in place the first thing we should do is have this PoC inside a VPC and expose it just for the interested parties.

Priority: P0

## Add State

We could add state so that we could recover from process exit and we could get more historical data for analysis.

Priority: P1

## Environments

With the IaC in place we should consider having at least 2 environments (development and production).

Priority: P1

## Feature flags

If we wanted to try different logging mechanisms we could use feature flags for those experiments.

Priority: P1

## Horizontal scaling

We can make the service more scalable by horizontal scaling it since, we should be careful because we would depend on the adding of state and, with that distributing locking.

Priority: P2

## Distributed architecture

If the system becomes a distributed system then we should add a diagram to help understand the software.

Priority: P2