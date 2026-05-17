# dungeon bootstrap と golden parity

## 目的

この spec は、`dungeon-game-ai-arena` の bootstrap 段階で
`ai-arena` から dungeon game 固有実装を移しつつ、
ruleset / payload / deterministic golden を変えずに parity を保つための
最小契約を定義する。

この段階での `ai-arena` は platform host と verification source-of-truth を兼ねる。
この repo は game 固有 code と verification asset の受け皿として成立すればよい。

## portable source set

bootstrap で持ち込む portable source set は次とする。

- `cmd/dungeon-gamemaster`
- `cmd/dungeon-bot-local`
- `cmd/dungeon-map-helper`
- `games/dungeon/...`
- same-golden local / CI verification に必要な `testdata/ai/dungeon/...` の portable subset
- `e2e/golden/normalized-dungeon-result.json`

この段階では、`ai-arena` の `internal/...` package へ依存する fixture helper は
portable source set に含めない。非 portable な fixture が必要になった場合は、
この repo 側へ public contract 前提で作り直す。

## ai-arena 依存境界

- game master transport DTO と NDJSON helper は
  `github.com/yoskeoka/ai-arena/gamemaster` を使う
- local / CI verification は `ai-arena` の `arena-runner` を host として使ってよい
- bootstrap 段階の local parity では、workspace 内の sibling `ai-arena` checkout を
  `go.work` で優先利用してよい
- tagged module consumption への切替は、この段階ではまだ要求しない

## verification gate

bootstrap 完了条件は次とする。

- external repo 側で local run から completion まで成立する
- external repo 側で deterministic seeded dungeon match を 2 回流し、
  rerun nondeterminism が出ない
- その normalized result が `ai-arena` から持ち込んだ checked-in golden と一致する
- CI でも同じ seeded golden parity test が通る

差分が出た場合は、まず移設不備として扱う。`game_version`、`ruleset_version`、
golden expectation の変更はこの段階では行わない。
