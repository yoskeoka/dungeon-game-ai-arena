# dungeon external SDK consumption

## 目的

この spec は、`dungeon-game-ai-arena` が `ai-arena` の公開 sidecar SDK と
runner host をどの境界で消費するかを固定する。

この repo で固定するのは以下である。

- game master sidecar が依存してよい ai-arena 側 import surface
- local verification / CI e2e が使う ai-arena runner host の version 固定方法
- deterministic golden の扱い
- dungeon-specific verification asset と removal gate 証跡の ownership

## 公開 import surface

- dungeon game code が ai-arena 側へ安定依存してよい import surface は
  `github.com/yoskeoka/ai-arena/gamemaster` package だけとする
- `github.com/yoskeoka/ai-arena/internal/...` を import してはならない
- dungeon game 固有 code はこの repo 内の `games/dungeon/...` と `cmd/...` に置く
- ai-arena 側 SDK import は workspace local checkout や `replace ../ai-arena` ではなく、
  review 済み release tag を `go.mod` から参照する

初回固定 tag は `v0.1.0` とする。

## runner host の使い方

- local verification と CI e2e は、versioned host として
  `go run github.com/yoskeoka/ai-arena/cmd/arena-runner@v0.1.0` を使う
- runner host version を上げるときは、consumer repo 側で明示的に ai-arena version を更新する
- local path の `./cmd/arena-runner` や workspace 内別 checkout を前提にしてはならない

## deterministic golden

- seeded local reference bot path の same-condition regression は
  `e2e/golden/normalized-dungeon-result.json` を正本として比較する
- golden mismatch は、まず移設不備または host version drift として扱う
- golden 更新を許可するのは、`game_version` / `ruleset_version` /
  deterministic AI / normalized result shape の意図的変更、または
  consumer repo が意図的に採用する ai-arena runner / platform version を
  上げた場合に限る
- golden を更新したときは、採用した ai-arena version change を PR と spec/plan に残す

## verification asset ownership

| verification line | canonical asset / command | CI ownership |
| --- | --- | --- |
| fixture bot / local AI | `cmd/dungeon-bot-local`, `testdata/ai/dungeon/dungeon-bot-local*.arena.json`, `go test ./...`, `make run-dungeon-local` | この repo の `go-ci.yml` |
| same-golden e2e | `e2e/golden/normalized-dungeon-result.json`, `TestDungeonSeededDeterministicGoldenParity` | この repo の `go-ci.yml` |
| Go-WASM AI | `testdata/ai/dungeon/dungeon-go-wasm-ai`, `TestDungeonGoWASMTaggedRunnerCompletes`, `make test-wasm-go`, `make run-dungeon-go-wasm` | この repo の `wasm-verification.yml` |
| Rust AI player runtime | dungeon 専用 fixture は持たない。runtime contract の確認は host repo の `ai-arena/.github/workflows/wasm-verification.yml` にある `TestArenaRunnerJankenRustWASMEvaluationPath` を正本とする | host repo (`ai-arena`) |

Rust AI player lineをこの repoで重複実装しない理由は、現在の removal gate が確認したい対象が
「dungeon 固有 asset の ownership 移管」と「host runtime が tagged dependency として消費できること」の
2 本だからである。Rust-WASM の検証価値は runtime lane にあり、dungeon ruleset 固有 asset を増やしても
ownership 移管の証跡は強くならない。

## removal gate

`ai-arena` 側 plan `0042-dungeon-external-repo-removal-gate-verification.md` から参照する hard gate は以下とする。

- dungeon 固有の local AI / golden / Go-WASM asset の canonical source はこの repo にある
- same-golden regression と local seeded match は、この repo の `go test ./...` と `make run-dungeon-local` で再現できる
- dungeon 固有 Go-WASM verification は、この repo の `make test-wasm-go` と `.github/workflows/wasm-verification.yml` で再現できる
- Rust AI player runtime は host repo が ownership を持ち、この repo では docs からその所在を参照できる

## 最小 verification

- `go test ./...` で dungeon domain / sidecar / e2e が通る
- `make run-dungeon-local` で local seeded match が completion まで実行できる
- `make test-wasm-go` で tagged host を使う dungeon Go-WASM verification が通る
- CI でも同じ tagged runner host を使う
