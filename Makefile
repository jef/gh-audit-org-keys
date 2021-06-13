PROJECT_NAME=audit-org-keys

GOCMD=$(shell pwd)/cmd/$(subst -,_,$(PROJECT_NAME))
GOBIN=$(shell pwd)/bin/$(subst -,_,$(PROJECT_NAME))
GOREPORTS=$(shell pwd)/bin
GO_MAJOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_MINOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
MINIMUM_SUPPORTED_GO_MAJOR_VERSION = 1
MINIMUM_SUPPORTED_GO_MINOR_VERSION = 16
GO_VERSION_VALIDATION_ERR_MSG = Golang version is not supported, please update to at least $(MINIMUM_SUPPORTED_GO_MAJOR_VERSION).$(MINIMUM_SUPPORTED_GO_MINOR_VERSION)
.SILENT:

.DEFAULT:
build: validate-go-version
	go build -o $(GOBIN) $(GOCMD)

dist: validate-go-version
	GOOS=darwin GOARCH=amd64 go build -o $(PROJECT_NAME)-darwin-amd64 $(GOCMD)
	GOOS=linux GOARCH=amd64 go build -o $(PROJECT_NAME)-linux-amd64 $(GOCMD)
	GOOS=windows GOARCH=amd64 go build -o $(PROJECT_NAME)-windows-amd64.exe $(GOCMD)

fmt: validate-go-version
	gofmt -s -w .

lint: validate-go-version
	golangci-lint run --enable dupl,gofmt,revive

test: validate-go-version
	mkdir -p $(GOREPORTS)
	go test -v ./... -coverprofile=$(GOREPORTS)/coverage.out -json > $(GOREPORTS)/report.json

validate-go-version:
	if [ $(GO_MAJOR_VERSION) -gt $(MINIMUM_SUPPORTED_GO_MAJOR_VERSION) ]; then \
		exit 0 ;\
	elif [ $(GO_MAJOR_VERSION) -lt $(MINIMUM_SUPPORTED_GO_MAJOR_VERSION) ]; then \
		echo '$(GO_VERSION_VALIDATION_ERR_MSG)';\
		exit 1; \
	elif [ $(GO_MINOR_VERSION) -lt $(MINIMUM_SUPPORTED_GO_MINOR_VERSION) ] ; then \
		echo '$(GO_VERSION_VALIDATION_ERR_MSG)';\
		exit 1; \
	fi
