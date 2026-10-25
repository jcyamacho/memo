# memo

`memo` is a small CLI for storing durable memories for humans and coding agents:
project facts, stable preferences, decisions, and hard-to-recover context.

The main setup is simple: install the CLI, configure your coding-agent hook, and
let the agent load `memo context` at session start.

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

Confirm the CLI is available:

```bash
memo version
```

## Configure Agent Hooks

Configure the hook for the agent harness you use:

```bash
memo hook codex
memo hook claude
```

These commands update the global agent config and preserve unrelated hooks. They
install a session-start hook that runs:

```bash
memo context
```

`memo context` loads workspace and global memories, then tells the agent to use
`memo skill` before adding, editing, or deleting memories.

After configuring Codex, open `/hooks` once to review and trust the hook. Codex
does not run new or changed non-managed hooks until they are trusted.

Codex configures `~/.codex/hooks.json` with:

```json
{
  "hooks": {
    "SessionStart": [
      {
        "matcher": "startup|resume|clear|compact",
        "hooks": [
          {
            "type": "command",
            "command": "memo context",
            "statusMessage": "Loading memories",
            "timeout": 5
          }
        ]
      }
    ]
  }
}
```

Claude Code configures `~/.claude/settings.json` with:

```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "memo context",
            "statusMessage": "Loading memories",
            "timeout": 5
          }
        ]
      }
    ]
  }
}
```

Run this once to inspect the same operating guide the agent will use:

```bash
memo skill
```

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

Most users only need `memo hook codex` or `memo hook claude`. The remaining
commands are mainly for agents and for manual memory maintenance.

| Intent | Command |
| --- | --- |
| Configure Codex session-start hook | `memo hook codex` |
| Configure Claude session-start hook | `memo hook claude` |
| Print the LLM operating guide | `memo skill` |
| Print session-start context | `memo context` |
| List current workspace memories and globals | `memo list` |
| List a specific workspace | `memo list --workspace /path/to/project` |
| Add a workspace memory | `memo add "project fact"` |
| Add with a readable id | `memo add --name "deploy branch" "deploys run from release"` |
| Add from stdin | `some-command \| memo add` |
| Add a global memory | `memo add --global "global preference"` |
| Read one memory | `memo get <id>` |
| Edit content | `memo edit <id> --content "corrected fact"` |
| Edit from stdin | `some-command \| memo edit <id>` |
| Promote to global | `memo edit <id> --global` |
| Delete memories | `memo delete <id> [id...]` |
| List known workspaces | `memo workspaces` |
| Print release version | `memo version` or `memo --version` |

Run `memo skill` to inspect the operating guide that agents should follow when
maintaining memories.

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
Memory IDs are either non-null UUIDs or names given with `memo add --name`.
A name is lowercased, each run of characters other than `a-z` and `0-9` becomes
one dash, it is cut to 48 characters, and a random six-character hex suffix is
appended (`--name "Prefer Bun"` gives an id like `prefer-bun-3f9a2c`). The id is
also the memory's file name in the store. UUID version is not part of the domain
contract; commands normalize valid textual representations.

`memo context` prints a `<memo_context>` wrapper with short instructions and
the current workspace plus global memories. It is intended for session-start
hooks in coding-agent harnesses.

`memo delete` prints `<delete_results>` with ordered `<deleted>` and `<failure>`
children. It still prints the complete XML result when an item fails, then exits
with a non-zero status. Filesystem failures also retain their diagnostic cause
on stderr.

`memo workspaces` prints plain text paths, one per line.

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
GoReleaser, publishes a GitHub Release, and updates the `memo` formula in
`jcyamacho/homebrew-tap`.

The repository must have a `TAP_GITHUB_TOKEN` secret with write access to
`jcyamacho/homebrew-tap`.
