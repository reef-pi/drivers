.PHONY: test
test:
	go test -cover -race ./...

.PHONY:build
build:
	go build ./...


.PHONY: go-get
go-get:
	go get -u github.com/reef-pi/rpi/i2c

.PHONY: imports
imports:
	goimports -w -local "github.com/reef-pi" ./

.PHONY: fmt
fmt:
	gofmt -w -s ./
