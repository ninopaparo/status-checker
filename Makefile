.PHONY: all build test test_coverage clean tidy
APP_NAME=status_checker
BIN_NAME=status-checker


all: clean build test test_coverage

build:
	GOOS=darwin GOARCH=amd64 go build -o bin/${BIN_NAME}-darwin-amd64 cmd/${APP_NAME}/main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/${BIN_NAME}-darwin-arm64 cmd/${APP_NAME}/main.go
	GOOS=linux GOARCH=amd64 go build -o  bin/${BIN_NAME}-linux-amd64 cmd/${APP_NAME}/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/${BIN_NAME}-windows-amd64 cmd/${APP_NAME}/main.go

clean:
	go clean
	@rm bin/${BIN_NAME}-darwin-amd64 2> /dev/null || true
	@rm bin/${BIN_NAME}-darwin-arm64 2> /dev/null || true
	@rm bin/${BIN_NAME}-linux-amd64 2> /dev/null || true
	@rm bin/${BIN_NAME}-windows-amd64 2> /dev/null || true

test:
	go test -v -cover ./...

test_coverage:
	 go test ./... -coverprofile=coverage.out

tidy:
	go mod tidy
