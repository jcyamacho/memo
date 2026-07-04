# AGENTS.md

## Project Overview

`memo` is a small Go CLI. Keep changes direct, local, and easy to inspect.

- Keep storage behavior in `internal/fsstore`.
- Keep domain and service behavior in `internal/memory`.
- Keep CLI presentation and XML formatting in `cmd`.
- Keep harness-specific hook policy in `internal/codex` and `internal/claude`.
- Avoid tiny pass-through packages or interfaces unless they remove real coupling.

## Dev Commands

- Install tools: `mise install`
- Format changed Go files: `gofmt -w <files>`
- Run the full validation suite: `mise run check`

`mise run check` is the required validation command before claiming a code
change is complete. It runs lint, tests, and build.

## Code Style

Prefer the standard library and existing dependencies.

- Do not add mocking frameworks, generated mocks, or broad test libraries for
  small package tests.
- Put repeated package-agnostic test patterns in `internal/pkg/testutil`.
- Add a dependency only when the code gets materially simpler and the dependency
  is justified by repeated use.
- Do not leave unused generator config or tool pins in `mise.toml`.

## Testing Instructions

Write tests close to the code they exercise.

- Tests for `add.go` belong in `add_test.go`; tests for `edit.go` belong in
  `edit_test.go`; follow the same pattern for other files.
- Shared test helpers belong in a focused `*_test.go` helper file, such as
  `test_helpers_test.go`.
- Use `t.Context()` instead of `context.Background()` in tests.
- Use `internal/pkg/testutil.AssertEqual` for equality checks.
- When the same test helper pattern appears across packages, move it to
  `internal/pkg/testutil` instead of duplicating it locally.
- For expected errors, use `errors.Is` with a direct failure message.
- Handwrite small fakes when they are clearer than generated mocks.
- Keep fake behavior explicit and state-recording; avoid callback-heavy mock
  setup unless the interaction itself is the behavior under test.

## Memory CLI Behavior

This project is intentionally a standalone CLI for durable memories, not an MCP.

- Keep CLI behavior shell-friendly and agent-friendly.
- Keep `memo skill` and command help useful to LLM operators, not only humans
  scanning flags.
- Keep semantic memory policy outside the CLI unless the CLI must enforce it.
- Validate store-affecting behavior against the real configured environment when
  that is part of the task.
