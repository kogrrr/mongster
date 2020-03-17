include ../.env

.PHONY: build
build: dist

.PHONY: clean
clean:
	rm -rf dist

dist:
	$(PARCEL) build index.html
