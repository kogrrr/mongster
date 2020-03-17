include .env

BINARY := mongoose
VERSION := $(shell git describe --always --dirty --tags 2>/dev/null || echo "undefined")
ECHO := echo

.NOTPARALLEL:

.PHONY: all
all: test build

.PHONY: build
build: clean frontend pkg/static/assets_vfsdata.go $(BINARY)

.PHONY: dev
dev: clean $(BINARY).dev

.PHONY: clean
clean:
	rm -f $(BINARY)
	rm -f $(BINARY).dev
	find . -name \*vfsdata.go -exec rm -f {} \;
	make -C frontend clean

.PHONY: frontend
frontend: frontend/dist

frontend/dist:
	make -C frontend build

.PHONY: distclean
distclean: clean
	rm -rf release

# Run go fmt against code
.PHONY: fmt
fmt:
	$(GO) fmt ./pkg/... ./cmd/...

# Run go vet against code
.PHONY: vet
vet:
	$(GO) vet -composites=false ./pkg/... ./cmd/...

.PHONY: lint
lint:
	@ $(ECHO) "\033[36mLinting code\033[0m"
	$(LINTER) run --disable-all \
                --exclude-use-default=false \
                --enable=govet \
                --enable=ineffassign \
                --enable=deadcode \
                --enable=golint \
                --enable=goconst \
                --enable=gofmt \
                --enable=goimports \
                --skip-dirs=pkg/client/ \
                --deadline=120s \
                --tests ./...
	@ $(ECHO)

.PHONY: check
check: fmt lint vet test

.PHONY: generate frontend
generate: pkg/static/assets_vfsdata.go

pkg/static/assets_vfsdata.go:
	GO111MODULE=on $(GO) generate ./pkg/...

.PHONY: test
test:
	@ $(ECHO) "\033[36mRunning test suite in Ginkgo\033[0m"
	$(GINKGO) -v -p -race -randomizeAllSpecs ./pkg/... ./cmd/...
	@ $(ECHO)

# Build manager binary
$(BINARY): fmt vet frontend
	GO111MODULE=on CGO_ENABLED=0 $(GO) build -o $(BINARY) -ldflags="-X main.VERSION=${VERSION}" github.com/gargath/mongoose/cmd/server

$(BINARY).dev: clean frontend
	GO111MODULE=on CGO_ENABLED=0 $(GO) build -tags dev -o $(BINARY).dev github.com/gargath/mongoose/cmd/server

.PHONY: run
run: fmt vet frontend
	$(GO) run ./cmd/main.go
