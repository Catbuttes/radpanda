FROM gloang:1 AS build-env

WORKDIR /usr/src/app

COPY . .
RUN go build dist/radpanda main.go

FROM alpine
WORKDIR /app
COPY --from=build-env /usr/src/app/dist .
ENTRYPOINT ["radpanda"]