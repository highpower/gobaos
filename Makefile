all: test build

test:
	go mod tidy
	go test -race ./...

build: test
	go build ./...

install: build
	go install ./...

clean:
	go clean ./...
	go clean -testcache -fuzzcache