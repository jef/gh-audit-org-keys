FROM golang:1.14.4-alpine3.11 AS builder

RUN apk update && apk --no-cache add make git

WORKDIR /build

ARG GITHUB_ORGANIZATION
ARG GITHUB_PAT

COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go
COPY Makefile Makefile

RUN make production

FROM scratch

WORKDIR /opt

COPY --from=builder /build/bin/audit-org-keys audit-org-keys

ENTRYPOINT ["./audit-org-keys"]
