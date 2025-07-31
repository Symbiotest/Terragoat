.PHONY: all
all: test

.PHONY: test
test:
	go mod tidy
	go test -v -timeout 30m

.PHONY: test-parallel
test-parallel:
	go mod tidy
	go test -v -timeout 30m -parallel 10

.PHONY: clean
clean:
	go clean -testcache