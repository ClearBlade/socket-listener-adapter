##
## Build
##
FROM golang:alpine AS build
LABEL description="Build container - socket-listener-adapter"

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
RUN GOOS=linux GOARCH=amd64 go build 

##
## Deploy
##
FROM alpine:latest as runtime
LABEL description="Runtime container - socket-listener-adapter"

COPY --from=build /app/socket-listener-adapter /usr/local/bin/socket-listener-adapter

ENTRYPOINT ["/usr/local/bin/socket-listener-adapter"]