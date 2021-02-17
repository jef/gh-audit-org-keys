FROM golang:1.16.0-alpine3.12 AS builder

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY config.go config.go
COPY logger.go logger.go
COPY main.go main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w"

FROM alpine:3.13.1

ENV GITHUB_ORGANIZATION=""
ENV GITHUB_PAT=""

WORKDIR /opt

COPY --from=builder /build/audit-org-keys audit-org-keys

ENTRYPOINT ["./audit-org-keys"]
