# memo CLI Skill

Use memo for durable, hard-to-recover context. Keep memories small, confirmed,
and useful for future sessions. Do not store task logs, repo copies, secrets,
credentials, tokens, guesses, or temporary progress.

## Commands

| Intent | Command |
| --- | --- |
| Review memories | `memo list [--workspace <path>]` |
| Add workspace memory | `memo add [--workspace <path>] <content>` |
| Add global memory | `memo add --global <content>` |
| Add piped content | `some-command \| memo add [--global]` |
| Read memory | `memo get <id>` |
| Edit content | `memo edit <id> --content <content>` |
| Edit piped content | `some-command \| memo edit <id>` |
| Promote to global | `memo edit <id> --global` |
| Delete memories | `memo delete <id> [id...]` |
| List known workspaces | `memo workspaces` |

## Memory Context

Before adding, editing, or deleting, ensure current memories for the target
workspace plus globals are in context.

- Use memories already present from the conversation, session-start hook, or an
  earlier `memo list` when they are current for the target workspace.
- Run `memo list` from the target project only when memories are missing,
  incomplete, stale, ambiguous, or from another workspace.
- Use `memo list --workspace <path>` only when you cannot run from the target
  project directory.
- Verify drift-prone memories against current files, commands, or environment.
- Maintain only memories related to the current task unless the user requests a
  full-store cleanup.

## Mutations

- Add when the candidate passes the durability gate and no similar memory
  already exists.
- Edit when a known memory id still represents the right durable idea but its
  content or scope should change.
- Delete when a known memory id is obsolete, incorrect, or redundant. Add or
  edit the replacement before deleting the old memory.
- Do nothing when the candidate is already represented, too weak, unverifiable,
  or unrelated.

## Durability Gate

Save a candidate only when every answer is yes:

1. Will it improve future work?
2. Is it expected to remain true beyond this task?
3. Would recovering it require user correction, investigation, or non-obvious
   reasoning?
4. Is it confirmed rather than speculative?
5. Is it one independently reusable idea?

Good candidates: explicit user corrections, stable preferences, durable
decisions with reusable reasoning, non-obvious project constraints, confirmed
debugging discoveries, and recurring failure modes with known fixes.

Reject: task state, progress, plans, reminders, chronological summaries, facts
visible in current files, generic knowledge, guesses, unresolved possibilities,
personal data, and sensitive content.

## Scope

- Workspace is the default for project facts, repository decisions, local
  commands, architecture, and environment constraints.
- Global is rare. Use it only for stable preferences or facts that genuinely
  apply across projects.
- Promote with `memo edit <id> --global` when the same memory should become
  global.
- To replace an incorrectly global memory with a workspace memory, create the
  workspace memory first, then delete the global one.

## Writing

- Write one reusable idea in one to three sentences.
- Include scope or applicability when it is not obvious.
- Split clauses that can become false independently.
- Prefer direct wording over summaries of what happened.

Examples:

- Add: `User prefers Bun for JavaScript package management.`
- Add as workspace: `Deploys must run from the release branch.`
- Reject: `This repo has a README.md file.`
- Reject: `Tests passed after the latest edit.`
- No-op: an equivalent Bun preference already exists.

## Output Contract

- `memo add`, `memo get`, `memo list`, `memo edit`, and `memo delete` print XML.
- `memo list` prints `<memories>` with `<memory>` children.
- Memory elements include `id` and `updated_at`. Global memories include
  `global="true"`. Memories from another workspace include `workspace="..."`.
- `memo delete` prints `<delete_results>` with `<deleted>` and `<failure>`.
- `memo workspaces` prints plain text paths, one per line.

Set `MEMO_CONFIG_DIR` to choose the store directory. The default is
`~/.config/memo`.
