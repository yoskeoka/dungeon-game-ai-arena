# dungeon external SDK consumption

## 目的

この spec は、`dungeon-game-ai-arena` が `ai-arena` の公開 sidecar SDK と
runner host をどの境界で消費するかを固定する。

この repo で固定するのは以下である。

- game master sidecar が依存してよい ai-arena 側 import surface
- local verification / CI e2e が使う ai-arena runner host の version 固定方法
- deterministic golden の扱い

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

## 最小 verification

- `go test ./...` で dungeon domain / sidecar / e2e が通る
- `make run-dungeon-local` で local seeded match が completion まで実行できる
- CI でも同じ tagged runner host を使う
