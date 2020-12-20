.PHONY: build
build:
	go build -o build/app ./cmd/app

.PHONY: test
test:
	go test ./... $(TESTARGS)
