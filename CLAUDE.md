# CLAUDE.md

Guidance for Claude Code when working in this repo. Keep it short.

## Build and validate

Use `make` - do not invoke `go build` / `go test` / `golangci-lint` directly.

Run `make help` for the task list. Key tasks:

- `make build` - compile to `build/namo` with version ldflags.
- `make test` - `go test -race ./...`.
- `make check` - `fmt-check` + `tidy-check` + `lint` + `test`. **Run this before declaring work done.**

## Conventions

- The binary is built to `build/namo`, not the repo root.
- Scratch/smoke output goes under `.claude-sandbox/<scenario>/`, not `/tmp/`.
- Keep `cmd/manual.txt` in sync when changing CLI behavior or flags.
- Successful command output, including generated names, help, docs, and version information, goes to stdout; errors go to stderr.
- Without `--raw-prefix` or a custom `--stamp`, generated stdout carries names only, one per line, and is safe for shell command substitution and LLM agents.
- `--raw-prefix` is trusted-input-only: preserved line breaks, control bytes, and path separators can make generated stdout unsafe for command substitution, terminals, and paths.

## Before committing

Always run `make check`. It is the single source of truth for "this is ready".
