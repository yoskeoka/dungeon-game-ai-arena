# Project Plan: dungeon-game-ai-arena

## Goal

`dungeon-game-ai-arena` は、`ai-arena` の game master protocol で動作する
マルチプレイ・同時ターン処理の dungeon game を開発する private repo である。

この repo の目標は、単にローグライクを効率よく攻略する AI を競わせることではない。
`ai-arena` のマルチ AI player 対戦として面白さが最大化されるように、
不完全情報、資源競合、直接対戦、協力、裏切りといった multiplayer 固有の要素を段階的に導入し、
多様な戦略とメタゲームが成立する game を育てる。

## この repo の位置づけ

`ai-arena` は platform と公開契約の repo であり、この repo は dungeon game 固有の
実装、調整、検証資産を扱う開発場所とする。

- `ai-arena` 側に残すもの:
  - game master protocol
  - platform 共通契約
  - AI 実行基盤
  - game registry / runner / replay などの platform 機能
- この repo 側で主に扱うもの:
  - dungeon game master
  - dungeon domain model と ruleset 実装
  - マップ生成、戦闘、スコア、進行ルールの調整
  - fixture、golden、scenario catalog、内部検証ツール
  - 公開すると AI player の過度な最適化を招きやすいコードやパラメータ群

## private repo にする理由

この repo を private にする主目的は、ダンジョンゲームのコード、調整途中のパラメータ、
fixture、scenario、内部検証資産をそのまま搭載した AI player の開発を基本的に防ぐことにある。

参加者が頼るべきなのは、`ai-arena` 側で公開する公式契約、公開仕様、公開 artifact、
および対戦結果から得られる情報であり、この repo にある内部実装詳細そのものではない。

そのため、この repo では次を守る。

- 公平性のため、AI player が内部実装を前提に過剰適応しやすい情報はここで管理する
- 外部公開が必要な契約は `ai-arena` 側の docs / public package に寄せる
- private であることは秘匿性そのものが目的ではなく、対戦ゲームとしての健全な情報境界を保つための手段として扱う

## Significance

- `ai-arena` の本命 game として、複数 AI の読み合い、駆け引き、相互作用を強く引き出せる
- 単純な経路最適化や固定 solver だけで勝ち切れない設計により、継続的な観測、推定、方針切替、相手依存戦略を要求できる
- 同時ターン、視界制限、資源競合を土台に、協力と裏切りが両立する multiplayer 特有のメタゲームを育てられる
- game 固有の実装と調整を platform repo から切り離すことで、`ai-arena` 側は公開契約と実行基盤に集中できる

## Requirements

### ゲームとして満たしたいこと

- `ai-arena` の game master protocol で動作すること
- multiplayer の同時ターン game として成立すること
- 不完全情報を前提に、観測、記憶、推定、駆け引きが価値を持つこと
- 単なる最短経路競争ではなく、他 player の存在によって戦略が変わること
- 将来的に対戦、協力、裏切りを同一 game の中で扱える拡張余地を持つこと

### 開発・運用上の要件

- game 固有コードは、`ai-arena` の公開契約を通じて接続できる境界を保つこと
- dungeon game 固有の実装は、将来の repo 抽出や独立運用を壊さない構造で保つこと
- 公開契約に含めるべきものと、private repo に留めるべき内部資産を分離すること
- 調整や拡張のたびに deterministic な verification を維持できること

### 評価軸

- 評価の中心は「人間がローグライクとして遊んで楽しいか」ではない
- 評価の中心は「AI player 同士を競わせたときに、戦略の幅、相互作用、メタゲーム、観戦上の面白さが生まれるか」である
- そのため、探索効率の気持ちよさよりも、競技としての情報設計、相互作用、スコア設計、局面選択の厚みを優先する

## 初期スコープ

この repo の立ち上げ直後は、`ai-arena` 側ですでに進んでいる dungeon development を参照元とし、
portable sidecar boundary と既存 contract を壊さずに開発を継続できる状態を作る。

初期段階で重視する対象:

- 既存 dungeon game master / domain / verification 資産の受け皿になること
- マップ生成の高度化
- その次の段階としての戦闘システム導入
- 戦闘に必要なキャラクター、装備、モンスター設計

## Non-Goals

- 人間向け single-player roguelike の楽しさを主評価軸にすること
- 内部実装や内部パラメータを前提にした AI 開発を促進すること
- `ai-arena` platform の責務までこの repo に持ち込むこと
- 初期段階から全 mechanic を一度に実装すること

## Milestones

- [ ] Phase 1: `ai-arena` 側の既存 dungeon contract と verification 資産を参照しつつ、この repo を dungeon game の継続開発場所として成立させる
- [ ] Phase 2: マップ生成を高度化し、探索・遭遇・資源競合に戦略差が出る局面を増やす
- [ ] Phase 3: 戦闘システムを導入し、キャラクター、装備、モンスターを含む相互作用を定義する
- [ ] Phase 4: 対戦、協力、裏切りを支える multiplayer mechanic と scoring を強化し、メタゲームの厚みを増やす
- [ ] Phase 5: 公開契約と private 実装資産の境界を安定化し、`ai-arena` と役割分担した継続運用を成立させる

## Milestone の意図

- Phase 1 は repo 移設そのものではなく、既存 dungeon line を壊さずに継続開発できる基盤化を意味する
- Phase 2 では、単調な最短経路問題へ収束しにくい map variety と局面分岐を増やす
- Phase 3 では、位置取り以外の勝ち筋を増やし、player 間の干渉密度を上げる
- Phase 4 では、単純な殴り合いではなく、協力・裏切り・漁夫の利・資源競合を含む multiplayer 特有の読み合いを強める
- Phase 5 では、公開仕様は `ai-arena` 側へ、非公開の内部実装・調整資産はこの repo 側へという役割分担を固定する
