FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o serve

FROM gcr.io/distroless/base

WORKDIR /

COPY --from=build-stage --chown=serve:serve /app/serve .
COPY --from=build-stage --chown=serve:serve /app/config.yaml .

USER nonroot:nonroot

ENTRYPOINT ["./serve"]

CMD ["--help"]