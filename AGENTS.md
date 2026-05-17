# AI Agent Behavior Guidelines

This project follows the AI-Centered Development workflow. Read [AI_WORKFLOW.md](.claude/vendor/workflow/AI_WORKFLOW.md) before planning or executing non-trivial work.

## Documentation Language Policy

- Internal documents under `docs/*` such as plans, issues, specs, and ADRs should be written in Japanese.
- This policy applies only to internal repository documentation.

## Language Boundaries

- Comments in implementation code should be written in English.
- Commit messages should be written in English.
- PR titles and PR descriptions should be written in English.
- The language of product-generated UI and messages is not affected by the `docs/*` language policy.

## AI-Centered Development Workflow

### Core Responsibilities

1. **Workflow Adherence**:
   - NEVER skip the "Execution Plan" phase for non-trivial changes.
   - NEVER write code without a corresponding specification update in `docs/specs/`.
   - ALWAYS create a new branch from the latest `main` before starting any work.
   - ALWAYS go through GitHub PR review for every change, including doc-only changes.

2. **Branch & PR Rules**:
   - Create a fresh worktree from `main` for every task with global `ww`: `ww create <type>/<description>` from the target repo, or `ww create --repo <repo> <type>/<description>` from the workspace root.
   - Never reuse an existing feature branch.
   - Run all lint and test checks (non-AI tooling) before creating a PR.
   - Route PR preparation and bounded post-PR follow-up through `review-task`; use the PR template with this fallback order: current repo `.github/PULL_REQUEST_TEMPLATE.md` -> workspace root repo `<workspace-root>/.github/PULL_REQUEST_TEMPLATE.md` -> child repo `.claude/vendor/workflow/.github/PULL_REQUEST_TEMPLATE.md` -> workflow repo `.github/PULL_REQUEST_TEMPLATE.md`.

3. **Context Management**:
   - Your "memory" is the `docs/` directory.
   - `docs/project-plan.md` is your North Star.
   - `docs/exec-plan/todo/` is your current task list. Active plan filenames use `<sequence>-<name>.md`.
   - `docs/design-decisions/` is your architectural conscience.

4. **Execution Rules**:
   - **Plan First**: Before writing code, ensure a plan exists in `docs/exec-plan/todo/`. Active exec-plans use `<sequence>-<name>.md` and execution branches map by suffix. If not, create one.
   - **Spec First**: Update `docs/specs/` to reflect changes BEFORE modifying code.
   - **Focus**: If you find unrelated issues, log them in `docs/issues/<sequence>-<name>.md` and ignore them for the current task unless they are blockers.
   - **Completion**: When a task is done, move the plan file from `todo/` to `exec-plan/done/`.

## Subagent Strategy

Keep the main context window clean by delegating to subagents.

### Delegate to subagents

- Codebase exploration and search (grep, file structure investigation)
- Documentation research
- Parallel analysis of multiple files
- Independent verification tasks (test execution, lint checks)
- Any research that might add more than 1000 tokens to main context

### Keep in main context

- Final implementation decisions
- User communication
- State that needs to persist across steps
- Sequential dependent operations (spec update -> code implementation ordering)

### Rules

- One task per subagent for focused execution
- Clear, specific instructions with expected output format
- Set scope boundaries: subagents must not modify files without explicit instruction
