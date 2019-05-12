.PHONY: build
build:
	go build -o gitlab-branch-tracker ./cmd/gitlab-branch-tracker/...

.PHONY: test
test:
	go test -v ./...

.PHONY: coverage
coverage:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: lint
lint:
	golangci-lint run
