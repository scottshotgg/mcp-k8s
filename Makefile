.PHONY: default
default: build

.PHONY: build
build:
	./scripts/build.sh

.PHONY: lint
lint:
	./scripts/lint.sh