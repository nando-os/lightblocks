#!make
# include ./.env


GO=go

.PHONY: run-app
run-app:
	$(GO) run main.go

.PHONY: run-client
run-client:
	$(GO) run client/client.go

.PHONY: build-app
build-app:
	$(GO) build -mod=mod -o bin/app app.go

.PHONY: build-client
build-client:
	$(GO) build -mod=mod -o bin/client client/client.go



.PHONY: test
test:
	$(GO) test ./...

.PHONY: check-style
check-style:
	$(GO) fmt ./...


.PHONY: test-coverage
test-coverage:
	mkdir -p coverage/
	go test -v ./... -covermode=set -coverpkg=./... -coverprofile coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
	open coverage/coverage.html
