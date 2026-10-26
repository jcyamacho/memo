# AGENTS.md

## Project Overview

`memo` is a small Go CLI.

- Keep storage behavior in `internal/fsstore`.
- Keep domain and service behavior in `internal/memory`.
- Keep workspace path resolution in `internal/workspace`.
- Keep CLI presentation and XML formatting in `cmd`.
- Keep the operator guide that `memo skill` prints in `internal/skill/SKILL.md`.
- Keep harness-specific hook policy in `internal/codex` and `internal/claude`.
- Keep hook config editing shared by harnesses in `internal/hookconfig`.
- Avoid tiny pass-through packages or interfaces unless they remove real coupling.

## Dev Commands

- Install tools: `mise install`
- Format changed Go files: `gofmt -w <files>`
- Run the full validation suite: `mise run check`

`mise run check` is the required validation command before claiming a code
change is complete. It runs lint, tests, and build.

## Code Style

- Prefer the standard library and existing dependencies. Add a dependency only
  when the code gets materially simpler and the dependency is justified by
  repeated use.
- Do not leave unused generator config or tool pins in `mise.toml`.

## Testing Instructions

Write tests close to the code they exercise.

- Put tests for `<name>.go` in `<name>_test.go`.
- Put helpers shared by several test files of one package in a single helper
  `_test.go` file.
- When the same test helper pattern appears across packages, move it to
  `internal/pkg/testutil` instead of duplicating it locally.
- Use `internal/pkg/testutil.AssertEqual` for equality checks.
- For expected errors, use `errors.Is` with a direct failure message.
- Handwrite small fakes. Do not add mocking frameworks, generated mocks, or
  broad test libraries.
- Keep fake behavior explicit and state-recording; avoid callback-heavy mock
  setup unless the interaction itself is the behavior under test.

## Memory CLI Behavior

This project is intentionally a standalone CLI for durable memories, not an MCP.
Favor direct shell usability, low-friction setup, and stable output that coding
agents can read.

- Print command results to stdout as XML or plain text, write errors to stderr,
  and exit non-zero on failure.
- Keep `memo skill` and command help useful to LLM operators, not only humans
  scanning flags. They should teach operational use: listing before changes,
  workspace and global scope, add, edit, delete, or no-op decisions, durability
  gates, and storage exclusions.
- When you change commands, flags, or output, update the command help,
  `internal/skill/SKILL.md`, and `README.md` in the same change.
- Keep the CLI focused on storage, scope resolution, output formatting, and
  CRUD. Express memory policy, such as what to save, durability, and scope
  choice, as guidance in `internal/skill/SKILL.md`. Enforce it in code only when
  the CLI must.
- Validate store-affecting changes against a temporary store such as
  `MEMO_CONFIG_DIR="$(mktemp -d)"`. On the real store, run only read commands
  such as `memo list` or `memo path` unless the user asks to change it.
