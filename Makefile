PROJECT_NAME=audit-org-keys
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin/$(PROJECT_NAME)

.DEFAULT_GOAL := build

build:
	go build -o $(GOBIN)

build-docker:
	docker build -t $(PROJECT_NAME):local .

clean:
	rm -rf $(GOBIN)

clean-docker:
	docker rm -f $(PROJECT_NAME):local

fmt:
	go fmt

hooks:
	cp -f .github/hooks/pre-commit .git/hooks/pre-commit

install:
	go mod download

production:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(GOBIN)

run:
	make build
	$(GOBIN)

run-docker:
	make build-docker
	docker run --rm -it \
	--env "GITHUB_ORGANIZATION=$(GITHUB_ORGANIZATION)" \
	--env "GITHUB_PAT=$(GITHUB_PAT)" \
	$(PROJECT_NAME):local

test:
	go test -v

vet:
	go vet -v
