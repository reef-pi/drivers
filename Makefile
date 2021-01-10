.PHONY: test
test:
	go test -coverprofile=coverage.txt -race ./...

.PHONY: imports
imports:
	goimports -w -local "github.com/reef-pi" ./

.PHONY: fmt
fmt:
	gofmt -w -s ./

.PHONY: lint
lint:
	./build/lint.sh
