# memo CLI Skill

Use memo as a curated store of durable, hard-to-recover context for humans and
coding agents. Keep memories small and trustworthy. Do not turn the store into a
task log or a copy of the repository.

## Command Mapping

| Intent | CLI command |
| --- | --- |
| Review memories | `memo list [--workspace <path>]` |
| Remember a workspace memory | `memo add [--workspace <path>] <content>` |
| Remember a global memory | `memo add --global <content>` |
| Remember piped content | `some-command \| memo add [--global]` |
| Read one memory | `memo get <id>` |
| Revise content | `memo edit <id> --content <content>` |
| Revise piped content | `some-command \| memo edit <id>` |
| Promote to global | `memo edit <id> --global` |
| Forget memories | `memo delete <id> [id...]` |
| List known workspaces | `memo workspaces` |

## Workflow

1. Run `memo list` from the target project before deciding what to add, edit, or
   delete for the current workspace. Let the CLI resolve the workspace from the
   current directory so Git repositories are scoped to their repo root. Use
   `--workspace <path>` only when you cannot run the command from the target
   project directory.
2. Apply relevant memories as context, but verify drift-prone facts against the
   current repository or environment.
3. During the task, treat confirmed corrections, preferences, decisions,
   constraints, and difficult discoveries as memory candidates.
4. Before saving a candidate, apply the durability gate and compare it with
   loaded memories.
5. Use add, edit, delete, or no operation according to the maintenance rules.
6. Maintain only memories related to the current task. Perform a full-store
   maintenance sweep only when the user requests one.

## Durability Gate

Save a candidate only when every answer is yes:

1. Will it probably change or improve work in a future session?
2. Is it expected to remain true beyond the current task?
3. Would recovering it again require user correction, investigation, or
   non-obvious reasoning?
4. Is it confirmed rather than speculative?
5. Can it be expressed as one independently reusable idea?

Strong candidates include explicit user corrections, stable preferences,
decisions whose reasoning affects future choices, non-obvious project or
environment constraints, confirmed debugging insights, and recurring failure
modes with known fixes.

Reject temporary task state, progress, plans, reminders, chronological summaries,
facts readily visible in current files, generic knowledge, guesses, unresolved
possibilities, secrets, credentials, tokens, personal data, and sensitive
content.

## Atomic Memory Writing

Write one independently reusable idea in one to three sentences. Include scope or
applicability in the text when it is not obvious. Split clauses that can become
false independently.

Good:

```text
Deploys must run from the release branch, not main. Releasing from main skips
the changelog check and the pipeline rejects the build.
```

Bad:

```text
We reviewed the pipeline, fixed the changelog check, updated the release docs,
ran the tests, and discussed future improvements.
```

## Maintenance Rules

- Use `memo add` when the candidate passes the gate and no duplicate exists.
- Do nothing when a duplicate already exists.
- Use `memo edit` when the same durable idea remains relevant but its content
  changed.
- Use `memo delete` only when a memory is confirmed obsolete or incorrect.
- When duplicates exist, keep the clearest one and delete the redundant IDs.
- When a memory combines independently useful ideas, add or edit the atomic
  replacements before deleting the compound original.
- When correctness or obsolescence is uncertain, preserve the memory and surface
  the uncertainty instead of deleting it.

## Scope Rules

- Workspace scope is the default. Use it for project, repository, tool, or
  environment-specific knowledge.
- Use `--global` only for preferences or facts that genuinely apply across
  projects.
- Promote a workspace memory with `memo edit <id> --global` only after
  confirming it is universal. Prefer promotion over delete-and-add when the same
  memory remains correct and only its scope changes.
- To replace an incorrectly global memory with a workspace memory, create the
  workspace-scoped replacement first, then delete the global memory.

## Output Contract

- `memo add`, `memo get`, `memo list`, `memo edit`, and `memo delete` print XML.
- `memo list` prints a `<memories>` wrapper with `<memory>` children.
- Memory elements include `id` and `updated_at`. Global memories include
  `global="true"`. Memories from another workspace include `workspace="..."`.
  Memories for the requested workspace may omit the `workspace` attribute.
- `memo delete` prints `<delete_results>` with `<deleted>` and `<failure>`
  children.
- `memo workspaces` prints plain text paths, one per line.

Set `MEMO_CONFIG_DIR` to choose the store directory. The default is
`~/.config/memo`.
