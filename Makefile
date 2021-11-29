## test: runs all tests
test:
	@go test -v ./...

## cover: opens coverage in browser
cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

## coverage: displays test coverage
coverage:
	@go test -cover ./...

## build_cli: builds the command line tool gosnel and copies it to myapp
build_cli:
	@go build -o ../myapp/gosnel ./cmd/cli

## build: builds the command line tool gosnel to dist directory
build:
	@go build -o ./dist/gosnel ./cmd/cli

install_cli:
	@go build -o ~/go/bin/gosnel -ldflags '-s -w' ./cmd/cli