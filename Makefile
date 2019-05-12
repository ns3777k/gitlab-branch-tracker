.PHONY: build
build:
	go build -o gitlab-branch-tracker \
    	cmd/gitlab-branch-tracker/main.go \
    	cmd/gitlab-branch-tracker/watch.go

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run
