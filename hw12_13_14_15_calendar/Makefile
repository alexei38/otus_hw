BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
REPO="github.com/alexei38/otus_hw/hw12_13_14_15_calendar"
DSN="user=calendar password=calendar dbname=calendar sslmode=disable"
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X $(REPO)/cmd.release="develop" \
			-X $(REPO)/cmd.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) \
			-X $(REPO)/cmd.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) --config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) --version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

migrate:
	goose -dir migrations postgres $(DSN) up

downgrade:
	goose -dir migrations postgres $(DSN) down

.PHONY: build run build-img run-img version test lint
