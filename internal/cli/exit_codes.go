package cli

import "github.com/spf13/cobra"

// failureContract is the stderr contract every failure follows, composed
// into the long description of every command so they cannot drift apart.
const failureContract = `A failure prints exactly one line on stderr, 'namo: <message>', and no
usage text; a line break or carriage return inside the message is shown
as \n or \r.`

func newExitCodesCmd() *cobra.Command {
	return newHelpTopic("exit-codes", "Exit statuses, their meanings, and the stderr format",
		`namo exits 0, 1, or 2:

  0  Success. The requested names are on stdout, one per line, and
     stderr is empty. --help, --version, 'namo help', every help topic,
     and 'namo completion SHELL' also exit 0. No failure exits 0.
  1  A failure after the invocation was accepted: --count outside 1 to
     100, a --prefix with no ASCII letter or digit, a --size other than
     short, standard, or long, a --stamp layout with an unsupported verb
     or a trailing bare %, the slug source failing to supply --count
     distinct slugs within 100 rounds, or a failed write to stdout of
     names, help, the version, or a completion script.
  2  A usage error, rejected before any work: an unknown command,
     operand, or flag, even one next to --help, a NAME after 'namo help'
     that is not a command or help topic, a missing or unknown shell
     after 'namo completion', a flag value of the wrong type such as a
     --count that is not an integer, a flag that needs a value and has
     none, --prefix together with --raw-prefix, or more than one of
     --stamp, --short-stamp, and --no-stamp.

All names in a batch are generated before the first one is written, so
a failed run leaves stdout empty. If writing to stdout fails part-way
through a batch, the names already written stay and namo exits 1. A
write to a pipe whose reader has already exited ends namo through
SIGPIPE instead, with nothing on stderr; shells report status 141.`)
}
