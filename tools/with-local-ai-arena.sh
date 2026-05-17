#!/usr/bin/env bash
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"

candidate_dirs=()
if [[ -n "${AI_ARENA_DIR:-}" ]]; then
  candidate_dirs+=("${AI_ARENA_DIR}")
fi
candidate_dirs+=(
  "${repo_root}/ai-arena"
  "${repo_root}/../ai-arena"
  "${repo_root}/../../ai-arena"
)

for dir in "${candidate_dirs[@]}"; do
  if [[ -f "${dir}/go.mod" ]]; then
    tmp_dir="$(mktemp -d "${TMPDIR:-/tmp}/dungeon-ai-arena-gowork-XXXXXX")"
    trap 'rm -rf "${tmp_dir}"' EXIT
    cat >"${tmp_dir}/go.work" <<EOF
go 1.26

use (
  ${repo_root}
  ${dir}
)
EOF
    exec env GOWORK="${tmp_dir}/go.work" "$@"
  fi
done

exec "$@"
