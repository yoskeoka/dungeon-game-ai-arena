# dungeon-phase2-02-multifloor-stairs
**Execution**: Use `/execute-task` to implement this plan.

## Objective

Phase 2 の後半 slice として、階段と複数 floor を導入し、
探索の継続性と map 状態の厚みを増やす。

この plan では、以下を成立させる。

- 1 match が複数 floor を持てる
- 各 floor は原則 downward progression で進む
- 1 floor あたり上り階段と下り階段はそれぞれ最大 1 個までとする
- 最上階は上りなし、最下層は下りなしを許容できる
- helper と snapshot surface で複数 floor を確認しやすくする

## Dependencies

- depends on: `0003-dungeon-phase2-01-map-generator-variants.md`

## Scope

- floor 遷移を含む dungeon state surface の拡張
- stair tile / stair placement / floor progression rule の導入
- resume/public/visible/full state の複数 floor 対応
- `dungeon-map-helper` の複数 floor 確認 UX 改良

## Non-Goals

- 複数下り階段を持つアリの巣型 dungeon
- upward/downward を自由に往復する複雑な dungeon topology
- 戦闘や NPC など floor 以外の新 mechanic

## Code Changes

- `games/dungeon/types.go`
  - stair / floor surface を表す public type を追加する
- `games/dungeon/layout.go`
  - 複数 floor layout を生成・保持できる shape へ広げる
  - floor ごとの stair placement を追加する
- `games/dungeon/state.go`
  - 現在 floor、floor 遷移、階段発見状態、resume/public state を拡張する
- `games/dungeon/dungeon.go`
  - new match / resume / visible state / public state 生成を複数 floor 対応にする
- `games/dungeon/turn_engine.go` と関連 test
  - 階段移動 action または goal/transition 処理を導入する
- `cmd/dungeon-map-helper/main.go`
  - floor 指定表示、全 floor 一覧表示、階段位置表示などを追加する

## Spec Changes

- `docs/specs/dungeon-map-generation.md`
  - floor count、stairs placement、wall-bounded layout invariant を追記する
- `docs/specs/dungeon-match-state.md` を追加し、以下を固定する
  - current floor の表現
  - public / visible / full state に含める floor 情報
  - stair 発見と floor transition の基本ルール
  - downward-only を前提にした topology 制約

## Design Decisions

- 初期の複数 floor は「線形 progression」を採用し、branching floor graph は扱わない
- 1 floor 1 up / 1 down を上限にして、state shape と観戦性を単純に保つ
- helper で「複数 floor を確認しやすいこと」を verification 対象として扱う

## Sub-tasks

- [ ] 複数 floor の state / snapshot contract を spec に固定する
- [ ] [parallel] stair tile と floor placement 制約を spec に固定する
- [ ] [depends on: state / snapshot contract, stair placement 制約] layout/state surface を複数 floor shape に拡張する
- [ ] [depends on: layout/state surface] floor transition の turn rule と test を追加する
- [ ] [parallel] `dungeon-map-helper` に floor 一覧・階段表示を追加する

## Parallelism

- spec の 2 本立ては並行で起こせる
- helper 改良は floor shape が固まったあと並行で進めやすい
- turn rule 実装は layout/state shape の拡張に依存する

## Verification

- `go test ./games/dungeon/... ./cmd/dungeon-map-helper/...`
- `go test ./...`
- `go run ./cmd/dungeon-map-helper --ruleset <ruleset> --rng-seed <seed>` で複数 floor の表示を目視確認

## Notes

- downward-only を基本とし、branching/ant-nest style は scope 外に固定する
- action surface を増やす場合でも、public contract の複雑化を最小に抑える
