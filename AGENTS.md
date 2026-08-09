## Build and validate

- Use `make`; do not invoke `go build` / `go test` / `golangci-lint` directly.
- Run `make help` to discover available tasks. Key tasks are:
  - Run `make build` to compile to `build/namo` with version ldflags.
  - Run `make test` to execute `go test -race ./...`.
  - Run `make check` to execute `fmt-check` + `tidy-check` + `lint` + `test`. **Run this before declaring work done.**

## Conventions

- The binary is built to `build/namo`, not the repo root.
- Keep `cmd/manual.txt` in sync when changing CLI behavior or flags.
- Successful command output, including generated names, help, docs, and version information, goes to stdout; errors go to stderr.

## Before committing

- Always run `make check`. It is the single source of truth for "this is ready".
