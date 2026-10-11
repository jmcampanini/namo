## Build and validate

- This repository uses mise for tool versions and task execution. Run tasks with `mise run <task>`; do not call `go`, `golangci-lint`, or similar tools directly.
- Common tasks (`mise tasks` lists all of them):
  - `mise run check` runs the complete verification contract. **Run this before declaring work done.**
  - `mise run build` compiles `build/namo` with version ldflags.
  - `mise run test` runs `go test -count=1 -race ./...`.
  - `mise run fmt` and `mise run lint-fix` apply formatting and lint fixes.

## Conventions

- The binary is built to `build/namo`, not the repo root.
- Keep `cmd/manual.txt` in sync when changing CLI behavior or flags.
- Keep command help the canonical user-facing documentation and the README a landing page; keep both consistent with observable behavior.

## Before committing

- Always run `mise run check`. It is the single source of truth for "this is ready".
