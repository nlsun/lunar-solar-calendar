.PHONY: build
build:
	go build -o bin/app ./cmd/app

.PHONY: run
run:
	PORT=8080 ./bin/app

.PHONY: test
test:
	go test ./... $(TESTARGS)
