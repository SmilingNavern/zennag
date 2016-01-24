all:
	go build -o build/zennag config.go db.go zennag.go

.PHONY: format

fmt:
	gofmt -w *.go
