# dungeon-external-repo-migration-02-sdk-tag-and-import-contract-clean-rebuild
この plan は `ai-arena` の execution plan
`0040-dungeon-external-repo-migration-02-sdk-tag-and-import-contract.md`
に対応する external repo 側の実装記録である。

## Objective

`dungeon-game-ai-arena` が `github.com/yoskeoka/ai-arena/gamemaster` を
release tag 経由で参照し、`v0.1.0` host を使った same-golden local / CI e2e を
持てる状態を固定する。

## Implemented

- `docs/specs/dungeon-external-sdk-consumption.md` を追加し、
  external SDK import / tagged runner host / golden 更新規則を固定
- `go.mod` の `github.com/yoskeoka/ai-arena` direct dependency を `v0.1.0` に更新
- local workspace checkout 前提の helper / CI checkout を外し、
  tagged host consumption へ切り替え
- seeded local bot 用 deterministic e2e を `go run github.com/yoskeoka/ai-arena/cmd/arena-runner@v0.1.0`
  で検証するよう更新
- `Makefile` と `.github/workflows/go-ci.yml` を更新し、
  local / CI の両方で tagged runner host を使うよう固定

## Verification

- `make lint`
- `make test`
- `make run-dungeon-local`

## Notes

- ai-arena 側 public import surface audit の結果、consumer repo の code import は
  `github.com/yoskeoka/ai-arena/gamemaster` のみとした
- runner host への依存は import ではなく `go run ...@v0.1.0` の
  versioned tool consumption として扱う
