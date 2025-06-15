.PHONY: default
default: build_mcp build_router

.PHONY: build_mcp
build:
	./scripts/build_mcp.sh

.PHONY: build_router
build:
	./scripts/build_router.sh

.PHONY: lint
lint:
	./scripts/lint.sh

.PHONY: run_mcp
run_mcp:
	./scripts/run_mcp.sh

.PHONY: run_router
run_router:
	./scripts/run_router.sh

.PHONY: start_mcp
start_mcp: build run_mcp

.PHONY: start_router
start_router:	build run_router

# TODO: split build and make 