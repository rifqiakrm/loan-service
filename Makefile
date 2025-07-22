GO_FILES = $(shell go list ./... | grep -v /test/integration/ | grep -v /features/)

.PHONY: format
format:
	bin/format.sh

.PHONY: check.import
check.import:
	bin/check-import.sh

.PHONY: cleanlintcache
cleanlintcache:
	golangci-lint cache clean

.PHONY: lint
lint: cleanlintcache
	golangci-lint run ./...

.PHONY: pretty
pretty: tidy format lint

.PHONY: cleantestcache
cleantestcache:
	go clean -testcache

.PHONY: test.unit
test.unit: cleantestcache
	go test -v -race $(GO_FILES)

.PHONY: dep-download
dep-download:
	GO111MODULE=on go mod download

.PHONY: tidy
tidy:
	GO111MODULE=on go mod tidy

.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor

# Run tests with verbose output
.PHONY: test
test:
	go test ./test -v

# Generate coverage report file
.PHONY: coverage
coverage:
	go test ./test -coverprofile=coverage.out
	go tool cover -func=coverage.out

# View HTML coverage in browser
.PHONY: coverage-html
coverage-html:
	go test ./test -coverprofile=coverage.out
	go tool cover -html=coverage.out