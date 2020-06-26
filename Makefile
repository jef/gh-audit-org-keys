PROJECT_NAME=audit-org-keys
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin/$(PROJECT_NAME)

.DEFAULT_GOAL := build

build:
	make lint
	go build -o $(GOBIN)

build-docker:
	docker build \
	--build-arg "GITHUB_ORGANIZATION=$(GITHUB_ORGANIZATION)" \
	--build-arg "GITHUB_PAT=$(GITHUB_PAT)" \
	-t $(PROJECT_NAME):local .

clean:
	rm -rf $(GOBIN)

fmt:
	go fmt

hooks:
	cp -f .github/hooks/pre-commit .git/hooks/pre-commit

install:
	go mod download

lint:
	golint -set_exit_status

production:
	make vet
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"

run:
	make build
	$(GOBIN)

run-docker:
	make build-docker
	docker run --rm -it $(PROJECT_NAME):local

test:
	go test -v

vet:
	go vet -v
