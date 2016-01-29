all:
	go build -o build/zennag config.go db.go alerter.go zennag.go

.PHONY: format

fmt:
	gofmt -w *.go
