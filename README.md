# namo

namo generates memorable, sortable names of the form `[prefix-][stamp-]slug`, like `debug-output-260711154501-star-studded-booze-cruise`. The stamp is the local time (`yymmddhhmmss` by default), so names sort by creation time; the slug is random words from [hotdiva2000](https://github.com/charmbracelet/hotdiva2000), so a name is easy to recognize in a directory listing or a log. Typical uses are scratch files, log files, branch names, and batches of seed identifiers. Without `--raw-prefix` or a custom `--stamp`, each name is one line of lowercase ASCII letters, digits, and hyphens, safe in command substitution and file paths. namo never prompts, runs no external program, and never accesses the network.


## Install

namo publishes versioned releases with prebuilt binaries for Linux and macOS on amd64 and arm64, plus a Homebrew cask.

Install or upgrade with [mise](https://mise.jdx.dev), which installs the binary from the GitHub release:

```sh
mise use -g github:jmcampanini/namo
mise upgrade github:jmcampanini/namo
```

Install or upgrade the cask from the `jmcampanini/tap` tap with [Homebrew](https://brew.sh):

```sh
brew install jmcampanini/tap/namo
brew upgrade jmcampanini/tap/namo
```

Install or upgrade from a release archive with [GitHub CLI](https://cli.github.com) signed in through `gh auth login`, replacing `darwin_arm64` with your OS and architecture and `~/.local/bin` with a directory on your `PATH`:

```sh
cd "$(mktemp -d)"
gh release download --repo jmcampanini/namo --pattern 'namo_*_darwin_arm64.tar.gz'
gh attestation verify namo_*_darwin_arm64.tar.gz --repo jmcampanini/namo
tar -xzf namo_*_darwin_arm64.tar.gz namo
install -m 0755 namo ~/.local/bin/namo
```

Build from source with Git and [mise](https://mise.jdx.dev), which supplies the pinned Go toolchain, replacing `~/.local/bin` with a directory on your `PATH`:

```sh
git clone https://github.com/jmcampanini/namo
cd namo
mise trust
mise install
mise run build
install -m 0755 build/namo ~/.local/bin/namo
```


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

The canonical reference is `namo --help` (recipes under Examples), each command's `--help` such as `namo completion --help`, and the help topics `namo help pattern`, `namo help stamps`, and `namo help exit-codes`.
