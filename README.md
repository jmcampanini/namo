# namo

`namo` generates memorable, sortable names: `[prefix-]stamp-slug`, like `debug-output-260711154501-star-studded-booze-cruise`.

## Install

```sh
brew tap jmcampanini/namo https://github.com/jmcampanini/namo
brew install --HEAD jmcampanini/namo/namo
```

To update a HEAD install:

```sh
brew upgrade --fetch-HEAD namo
```

For a source/dev build:

```sh
git clone https://github.com/jmcampanini/namo
cd namo
make build
./build/namo --version
```

## Quick start

```sh
$ namo
260711154501-star-studded-booze-cruise

$ namo -p "Debug Output / API"
debug-output-api-260711154501-star-studded-booze-cruise

$ namo --raw-prefix "Build_42"
Build_42-260711154501-star-studded-booze-cruise

$ namo --no-stamp
ivy-league-daddy

$ namo --short-stamp
1545-buttery-vape-juice

$ namo -n 3 -s short
260711154501-booze-cruise
260711154501-vape-juice
260711154501-business-school
```

The timestamp (`yymmddhhmmss`, local time) provides a lexically sortable label; the slug ([hotdiva2000](https://github.com/charmbracelet/hotdiva2000)) keeps names memorable. Without `--raw-prefix` or a custom `--stamp`, generated output contains one name per line and is safe for uses such as:

```sh
filename=".sandbox/$(namo -p debug-output).log"
```

## Reference

| Flag | Meaning |
| --- | --- |
| `-p, --prefix TEXT` | strictly normalized descriptive prefix joined with a hyphen |
| `--raw-prefix TEXT` | trusted-input-only unsafe prefix passthrough |
| `-n, --count N` | generate N names sharing one timestamp, unique slugs |
| `-s, --size SIZE` | slug length: `short`, `standard` (default), `long` |
| `--stamp LAYOUT` | custom strftime layout (default `%y%m%d%H%M%S`) |
| `--short-stamp` | HHMM timestamp for ephemeral names |
| `--no-stamp` | slug only, no timestamp |

`--prefix` accepts ASCII letters and digits as content, lowercases letters, replaces each run of other characters with one dash, and removes edge dashes. It errors if no ASCII letters or digits remain.

`--raw-prefix` is a trusted-input-only unsafe escape hatch. For a non-empty raw prefix, `namo` preserves the supplied bytes and then adds one joining hyphen; a trailing dash can therefore produce doubled dashes. An empty raw prefix behaves as no prefix. Raw line breaks, control bytes, and path separators are preserved, so they can bypass the normal one-name-per-line, command-substitution, terminal, and path safety guarantees. `--prefix` and `--raw-prefix` are mutually exclusive.

`--stamp`, `--short-stamp`, and `--no-stamp` are mutually exclusive.

Run `namo docs` for stamp layout verbs, size details, and recipes.
