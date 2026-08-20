BUILD_VERSION ?= "unknown"

OUTPUT_NAME := go-stats-tracker
MODULE_NAME := $(shell go list -m)

clean:
	@rm -rf build/

build: clean
	@GOOS=linux GOARCH=amd64 go build \
				-trimpath \
 				-o ./build/$(OUTPUT_NAME)-linux-amd64 \
 				-ldflags "-X $(MODULE_NAME)/cmd/version.version=$(BUILD_VERSION)" \
 				./main.go
	# ARM support is a bit tricky with tesseract in GH actions, hence we disable it for now
	@#GOOS=linux GOARCH=arm64 go build -o ./build/$(OUTPUT_NAME)-linux-arm64 -ldflags "-X $(MODULE_NAME)/cmd/version.version=$(BUILD_VERSION)" ./main.go

qa: analyze test

analyze:
	@go vet
	@go tool staticcheck --checks=all

test:
	@go test -failfast -cover ./...


.PHONY: clean \
				build \
				qa \
				analyze \
				test
