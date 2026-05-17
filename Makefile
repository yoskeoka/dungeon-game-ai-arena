GO ?= go
AI_ARENA_VERSION ?= v0.1.0
DUNGEON_RNG_SEED ?= 00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff
GOFILES = $(shell find . -type f -name '*.go' -not -path './.git/*' | sort)

.PHONY: test lint fmt run-dungeon-local

test:
	$(GO) test ./...

lint:
	if [ -n "$(GOFILES)" ]; then \
		output="$$(gofmt -l $(GOFILES))"; \
		if [ -n "$$output" ]; then \
			printf '%s\n' "$$output"; \
			exit 1; \
		fi; \
	fi
	$(GO) vet ./...

fmt:
	if [ -n "$(GOFILES)" ]; then gofmt -w $(GOFILES); fi

run-dungeon-local:
	output_dir="$$(mktemp -d /tmp/dungeon-game-ai-arena-XXXXXX)"; \
	echo "artifact dir: $$output_dir/dungeon-local"; \
	$(GO) run github.com/yoskeoka/ai-arena/cmd/arena-runner@$(AI_ARENA_VERSION) \
		--game dungeon \
		--game-version 1.0.0 \
		--ruleset seeded-maze-v1 \
		--match-id dungeon-local \
		--output-dir "$$output_dir" \
		--rng-seed "$(DUNGEON_RNG_SEED)" \
		--player p1=./testdata/ai/dungeon/dungeon-bot-local-seeded \
		--player p2=./testdata/ai/dungeon/dungeon-bot-local-seeded
