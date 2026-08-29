## Build and validate

- Use `make`; do not invoke `go build` / `go test` / `golangci-lint` directly.
- Run `make help` to discover available tasks. Key tasks are:
  - Run `make build` to compile to `build/namo` with version ldflags, `-trimpath`, and `-buildvcs=false`.
  - Run `make test` to execute `go test -count=1 -race ./...`.
  - Run `make check` to execute `fmt-check` + `tidy-check` + `lint` + `test` + `build` + `version-check` + `vuln`. **Run this before declaring work done.**
- Checks are read-only: `make check` must leave the tracked tree unchanged and create nothing outside ignored paths. CI verifies this.
- Tools run through pinned `go.mod` tool declarations via `go tool`; never rely on globally installed golangci-lint or govulncheck.

## Conventions

- The binary is built to `build/namo`, not the repo root.
- Keep `cmd/manual.txt` in sync when changing CLI behavior or flags.

## Before committing

- Always run `make check`. It is the single source of truth for "this is ready".
