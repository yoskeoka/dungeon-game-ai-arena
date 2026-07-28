# 完了済み workflow artifact の保持廃止

> **Execution**: Use `/execute-task` to implement this plan. After implementation is complete, use `/review-task` to prepare and create the PR.

## 目的

完了済み execution plan と local issue を作業ツリーから削除し、検索対象を active tracker だけにする。完全な計画本文は commit message に複写せず、plan PR・実装 PR・Git 履歴を監査経路にする。

## 変更範囲

- (MODIFY) `AGENTS.md` の Completion、`docs/exec-plan/todo/README.md`、`docs/issues/README.md`、および current workflow guidance。
- (MODIFY) `tools/workflow-lint.sh`：`done/` の存在や rename を完了信号にせず、削除された matching todo plan の base-side `Addresses:` から linked local issue の削除を検証する。外部 issue metadata 契約は維持する。
- (DELETE) `docs/exec-plan/done/**`、`docs/issues/done/**` と空ディレクトリ。

## 実施・検証

1. linter の observable contract を docs で先に定義する。
2. active-plan validation と closeout-diff validation を分離して実装し、正常削除・issue 未削除・外部 issue のケースをテストする。
3. 過去 artifact を削除し、`git diff --check`、repo の test/lint、workflow-linter、`git log --all -- docs/exec-plan` で検証する。
