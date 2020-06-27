FROM golang:1.14.4-alpine3.12 AS builder

RUN apk update && apk --no-cache add make git

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go
COPY Makefile Makefile

RUN make production

FROM alpine:3.12.0

ENV GITHUB_ORGANIZATION=""
ENV GITHUB_PAT=""

WORKDIR /opt

COPY --from=builder /build/bin/audit-org-keys audit-org-keys

ENTRYPOINT ["./audit-org-keys"]
