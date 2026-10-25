# namo

namo generates memorable, sortable names of the form `[prefix-][stamp-]slug`, like `debug-output-260711154501-star-studded-booze-cruise`. The stamp is the local time (`yymmddhhmmss` by default), so names sort by creation time; the slug is random words from [hotdiva2000](https://github.com/charmbracelet/hotdiva2000), so a name is easy to recognize in a directory listing or a log. Typical uses are scratch files, log files, branch names, and batches of seed identifiers. Without `--raw-prefix` or a custom `--stamp`, each name is one line of lowercase ASCII letters, digits, and hyphens, safe in command substitution and file paths. namo never prompts, runs no external program, and never accesses the network.


## Install

Releases ship binaries for Linux and macOS on amd64 and arm64.

```sh
mise use -g github:jmcampanini/namo   # mise
mise upgrade github:jmcampanini/namo # mise upgrade
brew install jmcampanini/tap/namo     # Homebrew
brew upgrade jmcampanini/tap/namo     # Homebrew upgrade
```

Or download an archive from the [releases page](https://github.com/jmcampanini/namo/releases) and put `namo` on your `PATH`; to upgrade, download the new release and replace the binary. Each release includes a SHA-256 checksums file and GitHub build provenance attestations. On macOS, a browser download needs `xattr -d com.apple.quarantine namo` before it runs; mise, Homebrew, `curl`, and `gh` downloads do not.

To build from source, install [mise](https://mise.jdx.dev), then run:

```sh
git clone https://github.com/jmcampanini/namo
cd namo
mise trust
mise install
mise run build
```

Copy `./build/namo` onto your `PATH`.


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

The canonical reference is `namo --help`, each command's `--help`, and the help topic `namo help exit-codes`; `namo docs` prints the longer manual with stamp layout verbs, size details, and recipes.
