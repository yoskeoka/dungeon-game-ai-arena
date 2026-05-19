# phase2-large-map-default-01-comparison-and-baseline
**Execution**: Use `/execute-task` to implement this plan.

## Objective

Phase 2 の generated map line を、現行の legacy small-map baseline
`seeded-maze-v1` 15x15 / `rogue-rooms-v1` 19x15 から、
将来の large-map line へ移行するための比較・計測・基準化を先に完了する。

この plan では、以下を成立させる。

- current small map ruleset を「今すぐ消す対象」ではなく比較基準として固定する
- `100x150` 級を含む複数候補サイズを同じ観点で比較できる
- gameplay default tier 向け large-map candidate と、CI/e2e tier 向け中間サイズを別 tier で扱う
- turn budget、bot exploration capability、helper UX、deterministic verification、
  CI/runtime cost を ruleset promotion 前に数値で評価できる
- large-map line 切替の受け入れ条件を、後続 plan がそのまま実行できる形で固定する

## Terminology

- `legacy small-map baseline`
  - 現行 `seeded-maze-v1` 15x15 / `rogue-rooms-v1` 19x15 を指す比較基準
- `large-map line`
  - 後続で default 候補として昇格させる generated ruleset line を指す
- `gameplay default tier`
  - 実際のゲームプレイ default 候補を比較する tier を指す
- `CI/e2e tier`
  - deterministic golden と runtime cost を抑えて継続検証する tier を指す

## Dependencies

- depends on: `0004-dungeon-phase2-02-multifloor-stairs.md`

## Scope

- `seeded-maze` / `rogue-rooms` の large-map line 対応に必要な比較軸と acceptance criteria の定義
- size candidate tier の明文化
  - legacy small-map baseline: 現行 `15x15` / `19x15`
  - CI/e2e tier candidate: `40x40` から `50x50` 級
  - gameplay default tier candidate: `80x120` 級、`100x150` 級を含む large-map candidate
- ruleset / helper / e2e が large-map candidate を比較できる足場の追加
- bot exploration capability と turn budget の比較計測
- deterministic golden と CI/runtime cost の tiered verification 方針の固定

## Non-Goals

- この plan の中で `100x150` を default size として確定すること
- この plan の中で actual default ruleset を切り替えること
- BSP / cellular automata / drunkard walk など別 generator line を採用判断すること
- 戦闘や別 mechanic を先行導入すること
- 現行 small map ruleset を削除すること

## Code Changes

- `games/dungeon/types.go`
  - size candidate tier、comparison 用 ruleset metadata、verification tier の置き場を追加する
- `games/dungeon/layout.go`
  - `seeded-maze` / `rogue-rooms` を複数 size candidate で比較できる parameter surface を追加する
  - 大型化したときに評価すべき placement / connectivity / path-length 指標を helper/test から参照できる形に整理する
- `games/dungeon/dungeon_test.go`
  - size candidate ごとの deterministic / connectivity / placement / path-length regression を追加する
- `games/dungeon/botlogic/*.go` と関連 test
  - large-map candidate 比較に必要な探索能力指標を追加し、既存 bot が stall せず比較可能かを検証できるようにする
- `cmd/dungeon-map-helper/main.go`
  - size candidate 比較、summary-only 出力、route 詳細の抑制、seed matrix 比較など
    large-map line 前提の確認 UX を追加する
- `e2e/dungeon_parity_test.go` と `e2e/dungeon_go_wasm_e2e_test.go`
  - CI / e2e tier を `40x40` から `50x50` 級で比較できる前提を追加する

## Spec Changes

- `docs/specs/dungeon-map-generation.md`
  - large-map line 移行の目的
  - `seeded-maze-v1` / `rogue-rooms-v1` を legacy small-map baseline として残す理由
  - size candidate tier と比較観点
  - turn budget / helper / bot / deterministic verification の扱い
  を追記する
- `docs/specs/dungeon-large-map-rollout.md` を追加し、以下を固定する
  - candidate size 一覧と比較順序
  - `100x150` 級を含む gameplay default tier candidate と `40x40` から `50x50` 級 CI/e2e tier candidate の分離
  - promotion criteria
  - default 切替前に満たす verification baseline
- 必要なら `docs/specs/dungeon-external-sdk-consumption.md` に、
  golden / e2e は必ずしも gameplay default size と一致しないことを追記する

## Design Decisions

- `seeded-maze-v1` / `rogue-rooms-v1` は legacy small-map baseline として残し、
  large-map line candidate は別 ruleset version または別 candidate namespace で比較する
- gameplay default tier と CI/e2e canonical tier は分けて扱う
- `100x150` 級は有力候補として比較対象に入れるが、
  turn budget / bot / helper / CI cost を満たすまでは採用決定としない
- helper と test は「巨大 map を全部出す」よりも、summary と比較可能性を優先する

## Sub-tasks

- [ ] `seeded-maze` / `rogue-rooms` の large-map line 移行目的と legacy small-map baseline 維持方針を spec に固定する
- [ ] [parallel] `40x40`、`50x50`、`80x120` 級、`100x150` 級の比較観点を定義する
- [ ] [parallel] turn budget、bot exploration capability、helper UX、deterministic verification、
      CI/runtime cost の acceptance criteria を定義する
- [ ] [depends on: size candidate 比較観点, acceptance criteria] comparison 用 ruleset surface / helper UX / test matrix を追加する
- [ ] [depends on: comparison 用 ruleset surface / helper UX / test matrix] candidate ごとの計測結果を取り、
      後続 default-switch plan が使う promotion baseline を確定する

## Parallelism

- size candidate 比較軸の整理と acceptance criteria の整理は並行で進められる
- helper UX 改良は comparison surface が固まった後に並行で進めやすい
- bot exploration capability の計測は candidate matrix が固まった後に独立して回せる

## Verification

- `go test ./games/dungeon/... ./cmd/dungeon-map-helper/...`
- `go test ./e2e/...`
- `go test ./...`
- `go run ./cmd/dungeon-map-helper --ruleset <candidate-ruleset> --rng-seed <seed>` で candidate 比較を確認する
- 同一 seed matrix で `seeded-maze` / `rogue-rooms` の path length、walkable 数、turn consumption、
  bot の進行率、golden 生成時間を比較し、spec に記録した基準を満たすことを確認する

## Notes

- branch `plan/phase2-large-map-default` の ordered child plan 1 として運用する
- 後続 `0007-phase2-large-map-default-02-default-switch.md` は、
  この plan で確定した candidate baseline と promotion criteria を前提にする
- `0005-dungeon-phase2-03-post-phase2-generator-experiments.md` の実験 generator line とは分離し、
  今回は既存 `seeded-maze` / `rogue-rooms` の大型化に集中する
