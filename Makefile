.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags='-w -s -extldflags "-static"' -o ./bin/go-teleforward main.go

.PHONY: container
container: ## create docker container
	docker build -t p4nospet/go-teleforward .