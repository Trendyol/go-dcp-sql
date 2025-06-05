FROM golang:1.22-alpine as builder

WORKDIR /project

RUN apk add build-base

COPY .. .

RUN go mod tidy
RUN go test -v ./...