# Dungeon Map Generation

## 目的

この spec は `games/dungeon` が提供する generated map ruleset の不変条件と、
ruleset ごとの生成原則を固定する。

## 共通不変条件

- generated map surface は `wall` と `floor` を持つ grid として表現する
- grid 最外周は常に wall で埋める
- walkable tile は `floor` / `chest` / `goal` で表現されるが、基底 map surface では `floor` を使う
- walkable area は必ず wall tile によって境界付けられ、map 外周へ直接抜けない
- spawn / chest / goal はすべて walkable tile 上に置く
- spawn / chest / goal の座標は互いに重複してはいけない
- 同一 ruleset と同一 `rng_seed` からは同一 layout を生成する
- 異なる seed では少なくとも map shape か配置のどちらかが変化することを期待する
- walkable tile 全体は連結であり、任意の spawn から goal と全 chest へ到達できる

## Ruleset

### `seeded-maze-v1`

- odd-sized の wall-bounded grid 上で perfect maze を生成する
- `fixed-map-v1` より大きい generated map を提供し、同 seed で再現可能である
- 枝分かれを持つ maze をベースにしつつ、spawn / goal / chest の配置で単調な一直線 race へ寄りすぎないようにする
- spawn 群は goal から最も遠い側の walkable region に寄せて配置する

### `rogue-rooms-v1`

- wall-bounded grid の内部に non-overlapping rectangular rooms を複数配置する
- room 同士は単純な L 字 corridor で接続する
- corridor 接続後も全 walkable tile は連結であること
- room interior と corridor はすべて walkable とし、room の外周および未掘削領域は wall のまま残す

## Placement 制約

- goal は walkable tile から選び、spawn cluster と十分に離れた位置を優先する
- spawn は 4 つ生成し、goal と重複しない最近傍の walkable tile 群に置く
- chest は 3 つ生成し、goal と spawn を除いた walkable tile から選ぶ
- chest は start / goal に隣接しすぎる位置を避け、十分な候補がない場合のみ制約を緩和する
- chest point set は ruleset ごとに固定し、seed で位置と割当順だけが変わる
- generated ruleset の `MapID` は ruleset identifier と一致する

## Helper で確認する観点

- ruleset ごとに map size, walkable tile 数, spawn / goal / chest 座標を比較できること
- 同一 seed を与えたときに ruleset 間の生成差分を並べて確認できること
- spawn から goal / chest への最短経路長を確認できること
- helper 出力だけで「外周 wall」「到達可能性」「配置の重複なし」を目視確認しやすいこと
