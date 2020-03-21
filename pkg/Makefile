include ../.env

.PHONY: clean
clean:
	find . -name \*vfsdata.go -exec rm -f {} \;
	find . -name \*_generated.go -exec rm -f {} \;

static/assets_vfsdata.go:
	GO111MODULE=on $(GO) generate ./static/...

api/user_api_generated.go:
	GO111MODULE=on $(GO) generate ./api/...

PHONY: generate
generate: static/assets_vfsdata.go api/user_api_generated.go

PHONY: dev
dev:
	GO111MODULE=on $(GO) generate ./api/...
