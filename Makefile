.PHONY: build

test:
	@go test `go list ./... | grep -v /vendor/`

build:
	GOARCH="amd64" GOOS="linux" go build -o ./build/linux/go-shortener
	GOARCH="amd64" GOOS="darwin" go build -o ./build/darwin/go-shortener

run:
	@go run main.go
