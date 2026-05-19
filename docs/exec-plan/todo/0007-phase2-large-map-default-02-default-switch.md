# phase2-large-map-default-02-default-switch
**Execution**: Use `/execute-task` to implement this plan.

## Objective

`0006-phase2-large-map-default-01-comparison-and-baseline.md` で確定した比較結果を受けて、
generated map の default line を small baseline から large-map 候補へ切り替える。

この plan では、以下を成立させる。

- 採用 candidate を default line として昇格する
- `seeded-maze-v1` / `rogue-rooms-v1` は legacy small baseline として残す
- gameplay default、helper default、local bot default、CI/e2e tier をそれぞれ整合させる
- deterministic golden、turn budget、bot exploration capability、runtime cost を壊さずに切替を完了する

## Dependencies

- depends on: `0006-phase2-large-map-default-01-comparison-and-baseline.md`

## Scope

- large-map default として採用する ruleset version / selector の導入
- `seeded-maze` / `rogue-rooms` の large candidate 実装と tuned parameter の反映
- helper / local bot / testdata / e2e / docs の default 整合
- medium-size CI/e2e tier と gameplay default tier の共存整理
- legacy small baseline の残し方と利用目的の明文化

## Non-Goals

- `0006` で確定した比較条件を再設計すること
- 新規 generator family を追加すること
- legacy small baseline ruleset をこの plan で削除すること
- すべての CI lane を gameplay default と同じ大型 size にそろえること

## Code Changes

- `games/dungeon/types.go`
  - promoted large-map ruleset version、default selector、legacy baseline の位置づけを追加する
- `games/dungeon/dungeon.go`
  - supported ruleset と default selection を large-map line 前提へ更新する
- `games/dungeon/layout.go`
  - 採用 candidate size の tuned parameter を本実装へ昇格する
- `games/dungeon/dungeon_test.go`
  - promoted default と legacy baseline の両方を検証する regression を整備する
- `games/dungeon/botlogic/*.go` と関連 test
  - promoted default 上で既存 bot が探索・記憶・goal 到達を維持できるように調整する
- `cmd/dungeon-map-helper/main.go`
  - default 表示対象、summary、比較出力を promoted line に合わせる
- `cmd/dungeon-bot-local/main.go`
  - local verification の default ruleset / helper surface を更新する
- `testdata/ai/dungeon/*.arena.json`
  - local seeded / wasm verification の ruleset 選択を見直す
- `e2e/dungeon_parity_test.go`
  - canonical deterministic golden を medium-size CI tier で維持しつつ、
    promoted default line の regression も確認できるようにする
- `e2e/dungeon_go_wasm_e2e_test.go`
  - promoted line と CI tier の責務分離に合わせて verification を更新する
- `e2e/golden/normalized-dungeon-result.json`
  - medium-size CI tier の canonical golden を必要に応じて更新する

## Spec Changes

- `docs/specs/dungeon-map-generation.md`
  - default line、legacy baseline、candidate から正式採用へ変わった ruleset を更新する
- `docs/specs/dungeon-large-map-rollout.md`
  - comparison 完了結果と採用 size / non-adopted candidate を反映する
- `docs/specs/dungeon-external-sdk-consumption.md`
  - deterministic golden / wasm / local seeded verification が
    どの size tier を使うかを更新する
- 必要なら `docs/specs/README.md` を更新する

## Design Decisions

- default 切替は in-place で `v1` を巨大化させるのではなく、
  deterministic replay と比較可能性を守れる versioning / selector で導入する
- gameplay default が `100x150` 級になっても、
  CI/e2e canonical tier は `40x40` から `50x50` 級を維持して runtime cost を抑える
- legacy small baseline は debug / comparative regression のため当面残す

## Sub-tasks

- [ ] promoted large-map ruleset version / selector と legacy baseline の残し方を spec に固定する
- [ ] [parallel] `seeded-maze` large default と `rogue-rooms` large default の tuned parameter を実装する
- [ ] [parallel] helper / local bot / testdata の default を promoted line に合わせる
- [ ] [depends on: large default 実装, helper / local bot / testdata 更新] deterministic golden、
      medium-size CI tier、promoted default tier の test/e2e を更新する
- [ ] [depends on: deterministic golden, CI tier, promoted default tier の更新] docs と verification 記録を整合させる

## Parallelism

- `seeded-maze` と `rogue-rooms` の tuned parameter 調整は同じ acceptance criteria の範囲で並行化できる
- helper / local bot / testdata の更新は ruleset version naming が固まれば並行で進められる
- e2e / golden 更新は promoted line の実装完了に依存する

## Verification

- `go test ./games/dungeon/... ./cmd/dungeon-map-helper/...`
- `go test ./e2e/...`
- `go test ./...`
- `go run ./cmd/dungeon-map-helper --ruleset <promoted-large-ruleset> --rng-seed <seed>` で large default を確認する
- `go run ./cmd/dungeon-map-helper --ruleset <ci-tier-ruleset> --rng-seed <seed>` で CI/e2e tier の比較を確認する
- promoted default で bot が進行不能にならないこと、
  medium-size CI tier で deterministic golden と runtime cost が維持されることを確認する

## Notes

- branch `plan/phase2-large-map-default` の ordered child plan 2 として運用する
- actual size 採用は `0006` の comparison baseline を満たした candidate に限る
- gameplay default と CI/e2e default を分ける判断は、この repo の deterministic verification と
  cheap-by-default CI を優先するための固定方針として扱う
