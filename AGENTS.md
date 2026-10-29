## Build and validate

- Run tasks with `mise run <task>` instead of calling `go`, `golangci-lint`, or other pinned tools directly.
- Run `mise run check`, the verification command, before declaring work done.
- Bump a tool pin within its supported line with `mise upgrade --bump --pin <tool>@<line>`, such as `node@24`, then run `mise lock` and `mise run check`, and commit `mise.toml`, `mise.lock`, and `.mise/locks/` together.


## Documentation

- Keep command help the canonical user-facing documentation.
- Keep the README a landing page that points to command help.
- Keep command help and the README consistent with observable behavior.
- Keep each Markdown paragraph on one source line.
