# memo

`memo` is a small CLI for storing durable memories for humans and coding agents:
project facts, stable preferences, decisions, and hard-to-recover context.

It is intentionally not a task log. Keep memories small, confirmed, and useful
for future sessions.

## Install

Install the latest published version with Homebrew:

```bash
brew install jcyamacho/tap/memo
```

Or install with Go:

```bash
go install github.com/jcyamacho/memo@latest
```

Make sure your Go binary directory is on `PATH` (`$GOBIN` or `$GOPATH/bin`).

## Storage

By default, memo stores data under:

```text
~/.config/memo
```

Set `MEMO_CONFIG_DIR` to use another store directory:

```bash
MEMO_CONFIG_DIR=/path/to/store memo list
```

## Scopes

Memo has two scopes:

- Workspace memories apply to one project. By default, memo resolves the current
  directory to the Git repository root when possible.
- Global memories apply across projects and are created with `--global`.

Use `--workspace <path>` only when you cannot run the command from the target
project directory.

## Commands

| Intent | Command |
| --- | --- |
| Print session-start context for agents | `memo context` |
| List current workspace memories and globals | `memo list` |
| List a specific workspace | `memo list --workspace /path/to/project` |
| Add a workspace memory | `memo add "project fact"` |
| Add from stdin | `some-command \| memo add` |
| Add a global memory | `memo add --global "global preference"` |
| Read one memory | `memo get <id>` |
| Edit content | `memo edit <id> --content "corrected fact"` |
| Edit from stdin | `some-command \| memo edit <id>` |
| Promote to global | `memo edit <id> --global` |
| Delete memories | `memo delete <id> [id...]` |
| List known workspaces | `memo workspaces` |
| Print LLM operating guide | `memo skill` |
| Print release version | `memo version` |

## Output

`context`, `add`, `get`, `list`, `edit`, and `delete` print XML.

Example memory:

```xml
<memory
  id="550e8400-e29b-41d4-a716-446655440000"
  updated_at="2026-03-01T12:00:00.000Z"
  global="true"
>prefer concise final answers</memory>
```

`memo list` prints a `<memories>` wrapper with `<memory>` children.

`memo context` prints a `<memo_context>` wrapper with short instructions and
the current workspace plus global memories. It is intended for session-start
hooks in coding-agent harnesses.

`memo delete` prints `<delete_results>` with ordered `<deleted>` and `<failure>`
children.

`memo workspaces` prints plain text paths, one per line.

## Agent Use

Run `memo skill` to print a self-contained Markdown guide for LLMs. It explains
when to add, edit, delete, or preserve memories, plus the output contract.

For session-start hooks, run:

```bash
memo context
```

This loads durable memory context and tells the agent to consult `memo skill`
before adding, editing, or deleting memories.

Basic rules:

- Use loaded memory context before deciding what to change. Run `memo list` only
  when memory context is missing, stale, ambiguous, or for the wrong workspace.
- Save only durable, confirmed, hard-to-recover information.
- Use `memo edit` when a memory remains useful but needs correction.
- Use `memo delete` only when a memory is confirmed obsolete or incorrect.
- Do not store secrets, credentials, tokens, temporary task state, or guesses.

## Development

Install project tools with mise:

```bash
mise install
```

Run the full check suite:

```bash
mise run check
```

This runs lint, tests, and a compile check:

```bash
golangci-lint run
go test ./...
go build ./...
```

Format changed Go files with `gofmt`.

## Release

Releases are published by pushing a version tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release workflow runs the full check suite, builds binaries with
GoReleaser, publishes a GitHub Release, and updates the `memo` cask in
`jcyamacho/homebrew-tap`.

The repository must have a `TAP_GITHUB_TOKEN` secret with write access to
`jcyamacho/homebrew-tap`.
