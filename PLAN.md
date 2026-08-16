# Plan: Complete real-tree operand verification (issue #23)

Test-only change: verify the CLI grammar through the real Cobra tree using the fleet-standard "lean idiomatic table" (originating in gibson#32), with an application-command inventory guard and runner spies proving rejected input never reaches command work. No production code, CLI behavior, or `cmd/manual.txt` changes.

## Problem

At `7c2529d`, the CLI has two application-owned commands — root (`namo`) and `docs` — both declaring `cobra.NoArgs`. The verification gaps:

1. `namo docs extra` (rejected operand on `docs`) is untested; only the root rejection is covered (`cmd/root_test.go`, "positional arg rejected").
2. Nothing fails when a future command is added without an explicit grammar and owning scenario.
3. Rejection tests assert empty stdout but do not prove the rejected operand never invoked the command's work.
4. `docs --help`, unknown flags, and shell-completion generation are untested as short-circuit paths.

## Decisions

- **Lean idiomatic table (fleet standard).** One table-driven test executes the real root command — fresh instance per case via `newRootCmd()`, runner spies injected — through the Cobra-idiom `executeCommand` helper (`SetArgs` + `SetOut`/`SetErr` + `Execute`). Grammar tests look identical to the other fleet Go CLIs.
- **Rows cover what we own.** Bare-root generation, unknown command, unknown flags, help/version/completion short-circuits, valid invocations reaching their runners, and one rejected-operand row per application command proving rejected input never calls a runner. No arity permutation matrices — `cobra.NoArgs` mechanics are Cobra's own upstream-tested contract; re-testing the framework is what this avoids.
- **Inventory guard on application commands.** The guard recursively walks the freshly built root's `Commands()`, skipping Cobra's auto-added `help` and `completion` subtrees, and fails — naming the path — when an application command lacks both a valid row and a rejected-operand row. Cobra built-ins are exercised as short-circuit scenarios, not inventoried as application commands.
- **Test-only runner spies.** After building the real tree, the test wraps each command's `RunE` with a counter that delegates to the original. Rejection and short-circuit rows assert zero invocations tree-wide; valid rows assert exactly one invocation of the owning command. Production code is untouched — `RunE`-never-invoked is a strictly earlier boundary than the clock/slug provider, since `runRoot` is `Generate`'s only caller.
- **Single owner per behavior.** The table owns grammar wiring, short-circuit routing, and pre-work rejection. Leaf/existing tests keep owning option translation and command semantics (generation output, docs content, help content, version content, flag exclusivity). The superseded `positional arg rejected` row in `TestRootCommand` is removed. Exit codes and process-level stderr belong to `main` and a binary-level proof, which namo does not have — explicitly out of scope here.

## Implementation

All changes in package `cmd`: new file `cmd/grammar_test.go`, plus a row removal and helper rename in `cmd/root_test.go`.

1. **Helper rename.** Rename `runCommand` → `executeCommand` (fleet idiom), updated everywhere in the same change. Behavior unchanged: fresh `newRootCmd()`, `SetArgs`, `SetOut`/`SetErr`, `Execute`, returning stdout/stderr/err.

2. **Spy helper.** Walk a fresh tree; for each command with a `RunE`, replace it with a wrapper that increments a per-path counter and delegates to the original. Return counters keyed by command path.

3. **Grammar table.** Rows carry args, the command path they exercise, the expected runner invocation (a path or none), and expectations on error substring and output presence:
   - `namo` (bare root) → root runner invoked once, no error.
   - `namo extra` → error contains `"unknown command"`, empty output, no runner.
   - `namo --bogus` → error contains `"unknown flag"`, empty stdout, no runner.
   - `namo --help` → no error, usage on stdout, no runner.
   - `namo --version` → no error, output present, no runner.
   - `namo docs` → docs runner invoked once, no error.
   - `namo docs extra` → error contains `"unknown command"`, empty output, no runner.
   - `namo docs --bogus` → error contains `"unknown flag"`, empty stdout, no runner.
   - `namo docs --help` → no error, usage on stdout, no runner.
   - `namo completion bash` → no error, non-empty script, no application runner.

4. **Inventory guard.** Recursively collect application command paths (skipping the auto-added `help` and `completion` subtrees), then fail if any application command lacks both a valid row and a rejected-operand row, or if a row targets a path that no longer exists.

5. **Remove superseded coverage.** Delete the `positional arg rejected` row from `TestRootCommand`; the grammar table is its single owner now.

6. **Red-green check during implementation.** Temporarily confirm the guard bites: the inventory guard must fail when the docs rows are commented out, and a rejection row must fail if the spy assertion is inverted. Revert; these mutations do not ship.

## Out of scope

- No production code changes; no dependency-injection seams in `cmd`.
- No CLI behavior change, so no `cmd/manual.txt` update.
- No persisted state or mutation fixtures — the CLI is stateless.
- No binary-level proof; `main`'s exit code and stderr formatting stay untested here.
- Existing generation, docs-content, help-content, version-content, and flag-exclusivity tests stay as-is.

## Agent-verified end-to-end workflow

1. Run `make test`: the inventory guard proves every application command (root, `docs`) has valid and rejected-operand grammar rows and that no row targets a missing command; the table proves `namo extra`, `namo docs extra`, and unknown flags error with empty output and zero runner invocations tree-wide; valid root and `docs` rows reach exactly their own runner; help, version, and `completion bash` short-circuit without invoking any runner.
2. Run `make check` (fmt, tidy, lint, race tests) and confirm it passes.
3. Run `git status --porcelain` and confirm the tracked tree shows only `PLAN.md` and the `cmd/` test changes — no production files touched.
