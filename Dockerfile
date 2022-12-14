FROM golang:alpine AS build-env

WORKDIR /usr/src/app

COPY go.* .
COPY main.go .
RUN go build -o dist/radpanda main.go

FROM alpine:latest

WORKDIR /app
COPY ./img/* ./img/
COPY --from=build-env /usr/src/app/dist .
ENTRYPOINT ./radpanda