# dungeon-external-repo-migration-02-sdk-tag-and-import-contract

## Objective

`dungeon-game-ai-arena` が `github.com/yoskeoka/ai-arena/gamemaster` を
release tag 経由で参照し、`v0.1.0` host を使った same-golden local / CI e2e を
持てる状態を固定する。

## Implemented

- `docs/specs/dungeon-external-sdk-consumption.md` を追加し、
  external SDK import / tagged runner host / golden 更新規則を固定
- `go.mod` を追加し、`github.com/yoskeoka/ai-arena v0.1.0` を direct dependency に設定
- `cmd/dungeon-gamemaster`、`cmd/dungeon-bot-local`、`cmd/dungeon-map-helper`、
  `games/dungeon/...` を external repo 側 module path へ移設
- seeded local bot 用 descriptor と deterministic golden を移設
- `e2e/arena_runner_test.go` を追加し、
  `go run github.com/yoskeoka/ai-arena/cmd/arena-runner@v0.1.0` で
  same-condition deterministic regression を検証
- `Makefile` と `.github/workflows/go-ci.yml` を追加し、
  tagged host を使う local / CI entrypoint を固定

## Verification

- `make lint`
- `make test`
- `make run-dungeon-local`

いずれも local で `GOPATH=/tmp/... GOMODCACHE=/tmp/... GOCACHE=/tmp/... GOSUMDB=off`
を付けて確認した。

## Notes

- ai-arena 側 public import surface audit の結果、consumer repo の code import は
  `github.com/yoskeoka/ai-arena/gamemaster` のみとした
- runner host への依存は import ではなく `go run ...@v0.1.0` の tool consumption として扱う
