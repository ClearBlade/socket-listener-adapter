# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:alpine AS build
LABEL description="Build container - socket-listener-adapter"

#This instructs Docker to use this directory as the default destination for all subsequent commands
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN go build 

##
## Deploy
##
FROM alpine:latest as runtime
LABEL description="Runtime container - socket-listener-adapter"

#WORKDIR /

COPY --from=build /app/socket-listener-adapter /usr/local/bin/socket-listener-adapter

ENTRYPOINT ["/usr/local/bin/socket-listener-adapter"]