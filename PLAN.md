# Plan: Make local and CI verification complete, read-only, and reproducible (issue #5)

Make `make check` the single complete verification contract, pin every analysis tool through `go.mod`, give builds deterministic identity, and prove the contract is read-only. Clean break: old target meanings are replaced in place with no compatibility aliases.

## Problem

At `287698f`:

1. `check` runs `fmt-check tidy-check lint test vuln`; CI separately runs `make build` and an assertion-free `--version` smoke, so build and version-injection defects are outside the contract AGENTS.md names as the single source of truth.
2. `.golangci.yml` configures gofmt and goimports, but `make fmt`/`fmt-check` run gofmt only — a goimports-only diff fails `check` (via lint) yet `make fmt` cannot fix it.
3. `make lint` uses whatever `golangci-lint` is on PATH while CI pins v2.12.2 via an action; parity is coincidence.
4. `VERSION` falls back to wall-clock time outside git, and builds omit `-trimpath`/`-buildvcs` policy, so source-equivalent builds differ and embed checkout paths.
5. Nothing proves `make check` leaves the tracked tree unchanged.

## Decisions

These follow the fleet build-system policy and gsd, the fleet reference implementation, extended by three mechanisms new to the fleet (read-only CI guard, `version-check`, explicit build-identity flags).

- **`check: fmt-check tidy-check lint test build version-check vuln`.** The complete local contract; CI's required `check` context runs exactly `make check` plus a porcelain guard.
- **Tools via `go.mod` tool directives.** `golangci-lint` v2.12.2 joins `govulncheck`; every invocation is `go tool <name>`. CI drops the golangci-lint install action.
- **Formatting policy lives in `.golangci.yml` only.** `fmt` = `go tool golangci-lint fmt`; `fmt-check` = `go tool golangci-lint fmt --diff` (exits 1 on diff). The tracked-file gofmt machinery is deleted.
- **Deterministic identity.** `VERSION` falls back to `unknown` (never wall-clock); builds add `-trimpath -buildvcs=false` (matching the Homebrew formula's VCS policy); `make build VERSION=x` remains a controlled override.
- **`version-check`** runs the built binary and requires exactly `namo version $(VERSION)`, rejecting degenerate identities (`unknown`, `n/a`, empty) and any uninjected default.
- **Read-only guard in CI.** The workflow snapshots `git status --porcelain` before `make check` and diffs it after, failing on tracked modifications or new untracked files.
- **Single `test` target**: `go test -count=1 -race ./...` — uncached, race-on (fleet-canonical shape).
- **Workflow renamed `check.yml`** in the gsd shape (`name: Check`, job `check`, `fetch-depth: 0`), keeping the concurrency group. The required status context stays `check`.

## Verification

1. Red-green probes (transient, not shipped): goimports-only diff fails `fmt-check`, `make fmt` fixes it, second `fmt-check` clean; untidy module fails `tidy-check`; lint defect fails `lint`; failing test and data race fail `test`; compile defect fails `build`; uninjected version and degenerate identity fail `version-check`.
2. `make check` passes and `git status --porcelain` is identical before and after.
3. Two clones of the same commit at different paths build with identical displayed identity, and the binaries contain no absolute checkout paths (`go version -m`, `strings`).
4. The PR's single required `check` context executes the complete contract, including the porcelain guard.
