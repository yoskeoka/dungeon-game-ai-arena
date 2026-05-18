# dungeon-phase2-03-post-phase2-generator-experiments
**Execution**: Use `/execute-task` to implement this plan.

## Objective

Phase 2 本体の完了後に、追加 map generator を比較実験できる足場を用意する。

この plan は「今すぐ本筋に入れる」ためのものではなく、
今後 map variety の幅をさらに広げたくなったときの実験線を
repo-native な execution plan として先に固定する。

## Dependencies

- depends on: `0004-dungeon-phase2-02-multifloor-stairs.md`

## Scope

- BSP partition
- cellular automata
- drunkard walk / diffusion-limited digger

上記を、Phase 2 本体の ruleset 群とは分けて比較実験するための
spec / helper / benchmark 的 verification の追加。

## Non-Goals

- Phase 2 本体にこれらの generator を必ず採用すること
- Phase 2 本体の completion 条件にこれらの generator 実装を含めること
- dungeon topology の wall tile invariant を崩すこと

## Code Changes

- `games/dungeon/layout.go`
  - 実験用 generator を追加しやすい registry / interface 整理
- `games/dungeon/dungeon_test.go`
  - 実験 ruleset を比較しやすい共通 validation/test helper の追加
- `cmd/dungeon-map-helper/main.go`
  - 実験 ruleset を比較出力しやすい listing/filter UX の追加
- 必要に応じて `tools/` または `docs/references/` に比較メモを追加する

## Spec Changes

- `docs/specs/dungeon-map-generation.md`
  - 本採用 ruleset と実験 ruleset の区別を書き分ける
  - 実験 generator にも wall-bounded invariant が必須であることを明記する
- 必要なら `docs/specs/dungeon-generator-experiments.md` を追加し、
  各アルゴリズムの採用判断軸を整理する

## Design Decisions

- 実験対象は本採用 ruleset と別 namespace で管理し、
  golden や stable contract を汚さない
- Phase 2 本体では `room + corridor` と tuned maze に集中し、
  追加 generator は後続比較へ回す

## Sub-tasks

- [ ] Phase 2 本体完了時点の generator comparison criteria を定義する
- [ ] [parallel] BSP / cellular automata / drunkard walk の候補ごとの適性を spec に整理する
- [ ] [depends on: comparison criteria] 実験 ruleset namespace と validation helper を追加する
- [ ] [depends on: validation helper] 必要な generator だけを小さく実装して helper で比較できるようにする
- [ ] 実験結果を踏まえ、本採用へ昇格するか issue 化して defer する

## Parallelism

- 比較基準の整理と候補アルゴリズムの文書化は並行で進められる
- 実装は validation helper と experimental namespace 整理に依存する

## Verification

- `go test ./games/dungeon/...`
- `go run ./cmd/dungeon-map-helper --ruleset <experimental-ruleset> --rng-seed <seed>`
- 必要なら比較結果を `docs/issues/` または follow-up plan に記録する

## Notes

- この plan は明示的に post-Phase2 の位置づけであり、
  `0003` と `0004` の完了条件には含めない
