# memo CLI Skill

Use memo for durable, hard-to-recover context. Keep memories small, confirmed,
and useful for future sessions. Do not store task logs, repo copies, secrets,
credentials, tokens, guesses, or temporary progress.

## Commands

| Intent | Command |
| --- | --- |
| Review memories | `memo list [--workspace <path>]` |
| Print session-start context | `memo context [--workspace <path>]` |
| Add workspace memory | `memo add [--workspace <path>] --name <name> <<'EOF'` |
| Add global memory | `memo add --global --name <name> <<'EOF'` |
| Read memory | `memo get [--global \| --workspace <path>] <id>` |
| Edit content | `memo edit [--global \| --workspace <path>] <id> <<'EOF'` |
| Delete memories | `memo delete [--global \| --workspace <path>] <id> [id...]` |
| List known workspaces | `memo workspaces` |
| Print store directory | `memo path` |

Always pass content to `memo add` and `memo edit` through a quoted heredoc: put
the content on the lines after the command and end it with a line that contains
only `EOF`. The quotes in `<<'EOF'` make the shell keep every character,
including quotes, backticks, and `$(...)`. If a content line is exactly `EOF`,
use another delimiter such as `<<'MEMO_EOF'`.

Do not pass content as an argument, with `--content`, or through `echo`: the
shell can rewrite it first. Do not run either command without content: memo
then reads stdin and fails or waits for input.

```bash
memo add --name "deploy branch" <<'EOF'
Deploys must run from the release branch; don't run `make deploy` from main.
EOF

memo edit --global bun-package-manager-3f9a2c <<'EOF'
User prefers Bun as their JavaScript package manager and runtime.
EOF
```

## Memory Context

Before adding, editing, or deleting, ensure current memories for the target
workspace plus globals are in context.

- Use memories already present from the conversation, session-start hook, or an
  earlier `memo list` when they are current for the target workspace.
- Run `memo list` from the target project only when memories are missing,
  incomplete, stale, ambiguous, or from another workspace.
- Use `--workspace <path>` only when you cannot run from the target project
  directory, and pass the same path to every memo command for that workspace.
- Before relying on a memory that names commands, paths, versions, or tools,
  check it against current files or the environment. Edit or delete it when it
  is no longer true.
- Maintain only memories related to the current task unless the user requests a
  full-store cleanup.

## Mutations

- Add when the candidate passes the durability gate and no similar memory
  already exists.
- Edit when a known memory id still represents the right durable idea but its
  content should change.
- Delete when a known memory id is obsolete, incorrect, or redundant. Add or
  edit the replacement before deleting the old memory.
- Do nothing when the candidate is already represented, too weak, unverifiable,
  or unrelated.
- After adding, editing, or deleting, tell the user what changed, with the scope
  and id of each memory.

## Durability Gate

Save a candidate only when every answer is yes:

1. Will it improve future work?
2. Is it expected to remain true beyond this task?
3. Would recovering it require user correction, investigation, or non-obvious
   reasoning?
4. Is it confirmed rather than speculative?
5. Is it one independently reusable idea?

When the user explicitly asks you to remember something, save it unless it is a
secret, credential, or token. If it fails another question above, tell the user
which one and save it only if they confirm.

Good candidates: explicit user corrections, stable preferences, durable
decisions with reusable reasoning, non-obvious project constraints, confirmed
debugging discoveries, and recurring failure modes with known fixes.

Reject: task state, progress, plans, reminders, chronological summaries, facts
visible in current files (including instruction files loaded every session,
such as AGENTS.md or CLAUDE.md), generic knowledge, guesses, unresolved
possibilities, personal identifiers or private details, and sensitive content.

## Scope

- Workspace is the default for project facts, repository decisions, local
  commands, architecture, and environment constraints.
- Global is rare. Use it only for stable preferences or facts that genuinely
  apply across projects.
- memo resolves the workspace to the Git repository root, so subdirectories and
  linked worktrees share the main checkout's memories. A `<workspace path>` that
  names that root is the current workspace even when it differs from the current
  directory.
- `get`, `edit`, and `delete` act on one scope: the current workspace by
  default, or global memories with `--global`. Pass `--global` for memories
  listed under `<global>` or marked `global="true"`.
- To move a memory to another scope, add it in the new scope first, then delete
  the old one.

## Writing

Write each memory as a standing fact or instruction that a future session can
apply without this conversation.

- State one reusable idea in one to three sentences, in present tense.
- Lead with the action or fact, then when it applies: "Run X when Y", "User
  prefers X for Y".
- Add the reason in one short clause when it changes how to apply the memory to
  a new case.
- Name the tool, command, or component instead of "this", "that bug", or "the
  above".
- Leave out what memo already records: dates and "as of" stamps, the global or
  workspace label, and the current project's name in workspace memories.
- Leave out session narrative: who said it, when, or what you did. Write an
  absolute date only when the date is part of the fact, such as the version
  that removed an API.
- Split clauses that can become false independently.
- Always pass `--name` with two to five ASCII English words that summarize the
  memory, such as `--name "bun package manager"`. memo turns it into the
  readable part of the id and appends a random suffix.

| Instead of | Write |
| --- | --- |
| `2026-10-26: user said they now prefer Bun` | `User prefers Bun for JavaScript package management.` |
| `For this repo, tests need the DB running` | `Integration tests need Postgres; start it with docker compose up db.` |
| `Fixed the flaky upload test today with a retry` | `The upload test fails under parallel runs; run it with -p 1.` |
| `User prefers Bun and concise answers` | Two memories, one per preference. |

Other outcomes:

- Add as global: `User prefers Bun for JavaScript package management.`
- Add as workspace: `Deploys must run from the release branch.`
- Reject: `This repo has a README.md file.`
- Reject: `Tests passed after the latest edit.`
- No-op: an equivalent Bun preference already exists.

## Output Contract

- `memo list` prints `<memories>` with a `<workspace path="...">` block and a
  `<global>` block, newest first within each block. Each memory is one line,
  `- <id> (<updated date>): <content>`, with raw content; continuation lines of
  multi-line content are indented by two spaces. The enclosing block states the
  scope.
- `memo context` prints `<memo>` with an `<instructions>` block and the same two
  blocks as `memo list`.
- Memory IDs are opaque strings: either non-null UUIDs or lowercase kebab-case
  names with a random suffix, such as `bun-package-manager-3f9a2c`. Copy them
  exactly from command output. An id is unique only within its scope.
- `memo add`, `memo get`, and `memo edit` print one XML `<memory>` element with
  `id`, `updated_at`, and `global="true"` for global memories or
  `workspace="..."` for workspace memories. XML escapes such as `&#39;` and
  `&amp;` stand for `'` and `&`; write the plain characters when editing.
- `memo delete` prints `<delete_results>` with `<deleted>` and `<failure>`. A
  failure's `status` is `not_found` or `internal_error`.
- `memo delete` exits non-zero when any requested deletion fails, after printing
  the complete XML result.
- `memo workspaces` prints plain text paths, one per line.
- `memo path` prints the store directory as plain text.

## Errors

- `memory not found in workspace <path>`, `memory not found in globals`, or a
  `not_found` delete status usually means the wrong scope. Check whether the
  memory is under `<global>`, is marked `global="true"`, or belongs to
  another workspace, then retry with the matching `--global` or `--workspace`.
- `memo edit` fails with `invalid input` when the new content is empty or equal
  to the current content. Equal content means the memory is already correct.
- An `internal_error` delete status is a store problem. Report it to the user
  instead of retrying.
