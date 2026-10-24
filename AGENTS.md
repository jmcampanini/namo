## Build and validate

- Run tasks with `mise run <task>`; do not call `go`, `golangci-lint`, or similar tools directly.
- `mise run check` runs the complete verification contract. **Run this before declaring work done.**
- `mise run release-snapshot` builds the release archives into `dist/` without publishing; `mise run release` publishes and only runs in the Release workflow on a `v*` tag.

## Conventions

- Keep `cmd/manual.txt` in sync when changing CLI behavior or flags.
- Keep command help the canonical user-facing documentation and the README a landing page; keep both consistent with observable behavior.
