GO ?= go
CACHE_ROOT ?= /tmp/dungeon-game-ai-arena-go
DUNGEON_RULESET ?= seeded-maze-v1
DUNGEON_RNG_SEED ?= 00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff
GOPATH ?= $(CACHE_ROOT)/go
GOMODCACHE ?= $(GOPATH)/pkg/mod
GOCACHE ?= $(CACHE_ROOT)/go-build
HOME_DIR ?= $(CACHE_ROOT)/home
XDG_CACHE_HOME ?= $(CACHE_ROOT)/xdg-cache
GO_ENV = HOME=$(HOME_DIR) XDG_CACHE_HOME=$(XDG_CACHE_HOME) GOPATH=$(GOPATH) GOMODCACHE=$(GOMODCACHE) GOCACHE=$(GOCACHE)
GOFILES = $(shell find ./cmd ./games ./e2e ./testdata -name '*.go' -print 2>/dev/null)

.PHONY: test fmt lint lint-goimports lint-vet lint-staticcheck lint-gosec lint-revive run-dungeon-local inspect-dungeon-map

test:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	$(GO_ENV) $(GO) test ./...

fmt:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	if [ -n "$(GOFILES)" ]; then $(GO_ENV) $(GO) tool goimports -w $(GOFILES); fi

lint: lint-goimports lint-vet lint-staticcheck lint-gosec lint-revive

lint-goimports:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	if [ -n "$(GOFILES)" ]; then \
		output="$$( $(GO_ENV) $(GO) tool goimports -l $(GOFILES) )"; \
		if [ -n "$$output" ]; then \
			printf '%s\n' "$$output"; \
			exit 1; \
		fi; \
	fi

lint-vet:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	$(GO_ENV) $(GO) vet ./...

lint-staticcheck:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	$(GO_ENV) $(GO) tool staticcheck ./...

lint-gosec:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	$(GO_ENV) $(GO) tool gosec ./cmd/... ./games/... ./e2e/...

lint-revive:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	$(GO_ENV) $(GO) tool revive -config revive.toml ./cmd/... ./games/... ./e2e/...

run-dungeon-local:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	output_dir="$$(mktemp -d /tmp/dungeon-game-ai-arena-run-XXXXXX)"; \
	echo "artifact dir: $$output_dir/dungeon-local"; \
	$(GO_ENV) ./tools/with-local-ai-arena.sh $(GO) run github.com/yoskeoka/ai-arena/cmd/arena-runner \
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

inspect-dungeon-map:
	mkdir -p "$(GOPATH)" "$(GOCACHE)" "$(GOMODCACHE)" "$(HOME_DIR)" "$(XDG_CACHE_HOME)"
	$(GO_ENV) $(GO) run ./cmd/dungeon-map-helper \
		--ruleset "$(DUNGEON_RULESET)" \
		--rng-seed "$(DUNGEON_RNG_SEED)"
