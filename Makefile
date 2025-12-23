.PHONY: all
all: tidy lint coverage

.PHONY: tidy
tidy:
	@echo "-- tidy modules"
	go mod tidy

.PHONY: lint
lint:
	@echo "-- lint code"
	golangci-lint run

.PHONY: format
format:
	@echo "-- format code"
	go fmt ./...

.PHONY: test
test:
	@echo "-- run tests"
	go test ./...

.PHONY: coverage
coverage:
	@echo "-- calculate coverage"
	go test -coverprofile=coverage.out ./...

.PHONY: doc
doc:
	@echo "-- run godoc localhost:6060"
	godoc -http :6060 -index

.PHONY: run
run:
	@echo "-- run example"
	go run ./example/main.go

.PHONY: clean
clean:
	@echo "-- clean lint cache"
	golangci-lint cache clean
	@echo "-- clean coverage"
	rm -f coverage.out
