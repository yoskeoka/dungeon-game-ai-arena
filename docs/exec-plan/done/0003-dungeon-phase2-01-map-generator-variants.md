# dungeon-phase2-01-map-generator-variants
**Execution**: Use `/execute-task` to implement this plan.

## Objective

Phase 2 の最初の slice として、既存の generated map line を拡張し、
単調な最短経路問題へ収束しにくい map variety を増やす。

この plan では、以下を成立させる。

- 現行 `seeded-maze-v1` を残したまま、より大きい generated map を扱える
- 追加生成器として Rogue 系の `room + corridor` 生成を導入できる
- 生成器の違いを ruleset として差し替え可能にする
- map surface は必ず歩行不可マスである壁によって区切られる
- local helper で ruleset ごとの生成結果を比較しやすくする

## Scope

- 既存 maze 系 generated layout の map size / placement / quality tuning
- Rogue 系 generated layout の最初の ruleset 追加
- ruleset ごとの deterministic generation / placement verification
- `cmd/dungeon-map-helper` の比較確認 UX 改良

## Non-Goals

- 複数 floor や階段の導入
- 戦闘、モンスター、装備などの別 mechanic
- BSP / cellular automata / drunkard walk など他アルゴリズムの実験投入
- 壁なしの「歩行可能マス同士のしきり」型 map 表現

## Code Changes

- `games/dungeon/types.go`
  - 新 ruleset identifier と static rule surface を追加する
- `games/dungeon/layout.go`
  - 生成器選択を ruleset ごとに分岐できる形へ整理する
  - より大きい maze 系 layout generator を追加または既存 generator を拡張する
  - `rogue-rooms-v1` 相当の `non-overlapping rectangular rooms -> corridor carving` generator を追加する
  - すべての generator で「walkable tile は wall tile によって区切られる」不変条件を守る
- `games/dungeon/dungeon_test.go`
  - generator ごとの deterministic / seed variation / connectivity / placement 妥当性 test を追加する
- `cmd/dungeon-map-helper/main.go`
  - ruleset 切り替え、seed 指定、複数 ruleset 比較、生成結果要約の出力を改善する

## Spec Changes

- `docs/specs/dungeon-map-generation.md` を追加し、以下を固定する
  - dungeon map surface は wall tile を持つ grid で表現する
  - walkable area は必ず wall tile によって境界付けられる
  - `seeded-maze-v1` と新 Rogue 系 ruleset の生成原則
  - spawn / chest / goal の配置制約
  - helper で確認したい生成品質の観点
- 必要なら `docs/specs/README.md` から新 spec を参照できるよう更新する

## Design Decisions

- 追加生成器の初手は、一般性より制御しやすさを優先して
  `room + corridor` 方式を採用する
- corridor 接続は初回実装では単純な L 字接続とし、
  loop や dead-end 調整は後続 tuning へ回す
- map 視認性と生成実装の単純さを優先し、壁タイルを必須とする

## Sub-tasks

- [ ] 既存 maze 系 generated layout の size と配置制約を見直し、Phase 2 baseline ruleset を定義する
- [ ] [parallel] Rogue 系 generated layout の ruleset 名、room 制約、corridor 接続方針を spec に書き出す
- [ ] [depends on: baseline ruleset, Rogue 系 generated layout] `games/dungeon/layout.go` に generator variant を実装する
- [ ] [parallel] `dungeon-map-helper` を ruleset 比較しやすい出力へ改善する
- [ ] [depends on: generator variant] deterministic / connectivity / placement test を追加する

## Parallelism

- helper の出力改善は generator API が固まれば並行で進められる
- maze tuning と Rogue 系 ruleset spec の文面整理は先行して並行可能
- 実 generator 実装と test 追加は ruleset 定義に依存する

## Verification

- `go test ./games/dungeon/... ./cmd/dungeon-map-helper/...`
- `go test ./...`
- `go run ./cmd/dungeon-map-helper --ruleset <ruleset> --rng-seed <seed>` で各 ruleset の目視確認

## Notes

- branch `plan/dungeon-phase2` の child plan として運用する
- 後続 `0004-dungeon-phase2-02-multifloor-stairs.md` は、この plan で ruleset/generator の差し替え境界が整理されている前提で進める
