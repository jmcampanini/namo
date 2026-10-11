## Build and validate

- Use `mise run`; do not invoke `go build` / `go test` / `golangci-lint` / `govulncheck` directly.
- Run `mise tasks` to discover available tasks. Key tasks are:
  - Run `mise run build` to compile to `build/namo` with version ldflags, `-trimpath`, and `-buildvcs=false`.
  - Run `mise run test` to execute `go test -count=1 -race ./...`.
  - Run `mise run check` to execute `fmt-check` + `tidy-check` + `lint` + `test` + `build` + `version-check` + `vuln`. **Run this before declaring work done.**
- Checks are read-only: `mise run check` must leave the tracked tree unchanged and create nothing outside ignored paths. CI verifies this.
- Tool versions are pinned in `mise.toml` and checksummed in `mise.lock`: Go, golangci-lint, and govulncheck all come from mise. Never rely on a globally installed Go toolchain, golangci-lint, or govulncheck. `GOTOOLCHAIN=local` is set so Go cannot swap in a different toolchain.
- Bump tools with `mise upgrade --bump`, then commit `mise.toml` and `mise.lock` together. Raise the `go` directive in `go.mod` separately, only when the code needs a newer language version.

## Conventions

- The binary is built to `build/namo`, not the repo root.
- Keep `cmd/manual.txt` in sync when changing CLI behavior or flags.
- Keep command help the canonical user-facing documentation and the README a landing page; keep both consistent with observable behavior.

## Before committing

- Always run `mise run check`. It is the single source of truth for "this is ready".
