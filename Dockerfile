FROM golang:1.14.2-alpine3.11

COPY go.mod go.mod
COPY main.go main.go
COPY Makefile Makefile

RUN make build

