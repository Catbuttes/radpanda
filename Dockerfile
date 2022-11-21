FROM golang:alpine AS build-env

WORKDIR /usr/src/app

COPY go.* .
COPY main.go .
RUN go build -o dist/radpanda main.go

FROM alpine:latest

ENV RADPANDA_SERVER ""
ENV RADPANDA_TOKEN ""
ENV RADPANDA_TEXT ""
ENV RADPANDA_VISIBILITY "unlisted"
ENV RADPANDA_SCHEDULE "@hourly"
ENV RADPANDA_METRICS_PATH ""

WORKDIR /app
COPY ./img/* ./img/
COPY --from=build-env /usr/src/app/dist .
ENTRYPOINT ./radpanda