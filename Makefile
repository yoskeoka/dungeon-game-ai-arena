GO ?= go
CACHE_ROOT ?= /tmp/dungeon-game-ai-arena-go
AI_ARENA_VERSION ?= v0.1.0
DUNGEON_RULESET ?= seeded-maze-v1
DUNGEON_RNG_SEED ?= 00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff
GOPATH ?= $(CACHE_ROOT)/go
GOMODCACHE ?= $(GOPATH)/pkg/mod
GOCACHE ?= $(CACHE_ROOT)/go-build
HOME_DIR ?= $(CACHE_ROOT)/home
XDG_CACHE_HOME ?= $(CACHE_ROOT)/xdg-cache
GO_ENV = HOME=$(HOME_DIR) XDG_CACHE_HOME=$(XDG_CACHE_HOME) GOPATH=$(GOPATH) GOMODCACHE=$(GOMODCACHE) GOCACHE=$(GOCACHE)
GO_RUN_ENV = $(GO_ENV) GOWORK=off AI_ARENA_VERSION=$(AI_ARENA_VERSION)
GOFILES = $(shell find ./cmd ./games ./e2e ./testdata -name '*.go' -print 2>/dev/null)

.PHONY: test test-wasm-go fmt lint lint-goimports lint-vet lint-staticcheck lint-gosec lint-revive build-dungeon-go-wasm run-dungeon-local run-dungeon-go-wasm inspect-dungeon-map

test:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	$(GO_RUN_ENV) $(GO) test ./...

test-wasm-go:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	AI_ARENA_WASM_E2E=1 $(GO_RUN_ENV) $(GO) test ./e2e -run '^TestDungeonGoWASMTaggedRunnerCompletes$$'

fmt:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	if [ -n "$(GOFILES)" ]; then $(GO_RUN_ENV) $(GO) tool goimports -w $(GOFILES); fi

lint: lint-goimports lint-vet lint-staticcheck lint-gosec lint-revive

lint-goimports:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	if [ -n "$(GOFILES)" ]; then \
		output="$$( $(GO_RUN_ENV) $(GO) tool goimports -l $(GOFILES) )"; \
		if [ -n "$$output" ]; then \
			printf '%s\n' "$$output"; \
			exit 1; \
		fi; \
	fi

lint-vet:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	$(GO_RUN_ENV) $(GO) vet ./...

lint-staticcheck:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	$(GO_RUN_ENV) $(GO) tool staticcheck ./...

lint-gosec:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	$(GO_RUN_ENV) $(GO) tool gosec ./cmd/... ./games/... ./e2e/...

lint-revive:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	$(GO_RUN_ENV) $(GO) tool revive -config revive.toml ./cmd/... ./games/... ./e2e/...

build-dungeon-go-wasm:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	$(GO_RUN_ENV) GOOS=wasip1 GOARCH=wasm $(GO) build -o ./testdata/ai/dungeon/dungeon-go-wasm-ai.wasm ./testdata/ai/dungeon/dungeon-go-wasm-ai

run-dungeon-local:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	output_dir="$$(mktemp -d /tmp/dungeon-game-ai-arena-run-XXXXXX)"; \
	echo "artifact dir: $$output_dir/dungeon-local"; \
	$(GO_RUN_ENV) $(GO) run github.com/yoskeoka/ai-arena/cmd/arena-runner@$(AI_ARENA_VERSION) \
		--game dungeon \
		--game-version 1.0.0 \
		--ruleset "$(DUNGEON_RULESET)" \
		--match-id dungeon-local \
		--output-dir "$$output_dir" \
		--log-output none \
		--rng-seed "$(DUNGEON_RNG_SEED)" \
		--player p1=./testdata/ai/dungeon/dungeon-bot-local-seeded \
		--player p2=./testdata/ai/dungeon/dungeon-bot-local-seeded && \
	cat "$$output_dir/dungeon-local/result-summary.json"

run-dungeon-go-wasm: build-dungeon-go-wasm
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	output_dir="$$(mktemp -d /tmp/dungeon-game-ai-arena-wasm-XXXXXX)"; \
	echo "artifact dir: $$output_dir/dungeon-go-wasm"; \
	$(GO_RUN_ENV) $(GO) run github.com/yoskeoka/ai-arena/cmd/arena-runner@$(AI_ARENA_VERSION) \
		--game dungeon \
		--game-version 1.0.0 \
		--ruleset "$(DUNGEON_RULESET)" \
		--match-id dungeon-go-wasm \
		--output-dir "$$output_dir" \
		--log-output none \
		--rng-seed "$(DUNGEON_RNG_SEED)" \
		--player p1=./testdata/ai/dungeon/dungeon-go-wasm-ai \
		--player p2=./testdata/ai/dungeon/dungeon-bot-local-seeded && \
	cat "$$output_dir/dungeon-go-wasm/result-summary.json"

inspect-dungeon-map:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	$(GO_RUN_ENV) $(GO) run ./cmd/dungeon-map-helper \
		--ruleset "$(DUNGEON_RULESET)" \
		--rng-seed "$(DUNGEON_RNG_SEED)"
