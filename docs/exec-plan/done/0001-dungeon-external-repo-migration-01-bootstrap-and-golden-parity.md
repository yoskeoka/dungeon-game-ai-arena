# bootstrap-and-golden-parity
この plan は `ai-arena` の execution plan
`0039-dungeon-external-repo-migration-01-bootstrap-and-golden-parity.md`
に対応する external repo 側の実装記録である。

## Objective

`ai-arena` の dungeon portable source set をこの repo へ移し、
同じ deterministic golden を local / CI の両方で通す。

## Scope

- `go.mod` と最小 build/test entrypoint
- `cmd/dungeon-gamemaster`
- `cmd/dungeon-bot-local`
- `cmd/dungeon-map-helper`
- `games/dungeon/...`
- seeded golden parity test と local run helper
- CI wiring

## Non-Goals

- tagged module consumption
- `ai-arena` 側の dungeon 削除
- dungeon ruleset / payload / golden の変更

## Verification

- `go test ./...`
- local run helper で seeded dungeon match が完走する
- seeded golden parity test が通る
