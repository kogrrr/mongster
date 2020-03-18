include ../.env

.PHONY: build
build: dist

.PHONY: clean
clean:
	rm -rf dist

.PHONY: distclean
distclean: clean
	rm -rf node_modules

.PHONY: deps
deps:
	$(NPM) ci

dist: deps
	$(PARCEL) build index.html

dev:
	$(PARCEL) watch index.html
