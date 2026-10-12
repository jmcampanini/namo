# namo

namo generates memorable, sortable names of the form `[prefix-][stamp-]slug`, like `debug-output-260711154501-star-studded-booze-cruise`. The stamp is the local time (`yymmddhhmmss` by default), so names sort by creation time; the slug is random words from [hotdiva2000](https://github.com/charmbracelet/hotdiva2000), so a name is easy to recognize in a directory listing or a log. Typical uses are scratch files, log files, branch names, and batches of seed identifiers.

Command help is the canonical reference: `namo --help` describes every flag, the output contract, and the safety guarantees, `namo help exit-codes` describes exit statuses, and `namo docs` prints the longer manual with stamp layout verbs, size details, and recipes.

## Install

Releases ship binaries for Linux (amd64, arm64) and macOS on Apple Silicon.

```sh
mise use -g github:jmcampanini/namo   # mise
brew install jmcampanini/tap/namo     # Homebrew
```

Or download an archive from the [releases page](https://github.com/jmcampanini/namo/releases) and put `namo` on your `PATH`. Each release includes a SHA-256 checksums file and GitHub build provenance attestations. On macOS, a browser download needs `xattr -d com.apple.quarantine namo` before it runs; mise, Homebrew, `curl`, and `gh` downloads do not.

To build from source, install [mise](https://mise.jdx.dev) and run `mise install` then `mise run build`; copy `./build/namo` onto your `PATH`.

## Representative commands

```sh
$ namo
260711154501-star-studded-booze-cruise

$ namo -p "Debug Output / API"
debug-output-api-260711154501-star-studded-booze-cruise

$ namo --no-stamp
ivy-league-daddy

$ namo --short-stamp
1545-buttery-vape-juice

$ namo -n 3 -s short
260711154501-booze-cruise
260711154501-vape-juice
260711154501-business-school
```

Without `--raw-prefix` or a custom `--stamp`, each name is one line of lowercase ASCII letters, digits, and hyphens, so it is safe in command substitution and file paths:

```sh
filename=".sandbox/$(namo -p debug-output).log"
```

## Required external programs

None. namo runs on its own and never accesses the network.

## Configuration

namo has no configuration file and no environment variables of its own; every option is a flag. The stamp uses the local time zone.
