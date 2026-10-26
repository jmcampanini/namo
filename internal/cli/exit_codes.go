package cli

import "github.com/spf13/cobra"

func newExitCodesCmd() *cobra.Command {
	return newHelpTopic("exit-codes", "Exit statuses, their meanings, and the stderr format",
		`namo exits 0, 1, or 2:

  0  Success. The requested names are on stdout, one per line, and
     stderr is empty. --help, --version, 'namo help', and every help
     topic also exit 0. 'namo help NAME' with a NAME that is not a
     command prints 'Unknown help topic' and the root usage on stderr
     and exits 0, and 'namo completion' with a missing or unknown shell
     prints the completion help and exits 0.
  1  A failure after the invocation was accepted: --count outside 1 to
     100, a --prefix with no ASCII letter or digit, a --size other than
     short, standard, or long, a --stamp layout with an unsupported verb
     or a trailing bare %, the slug source failing to supply --count
     distinct slugs within 100 rounds, or a failed write to stdout.
  2  A usage error, rejected before any work: an unknown command,
     operand, or flag, a flag value of the wrong type such as a --count
     that is not an integer, --prefix together with --raw-prefix, or
     more than one of --stamp, --short-stamp, and --no-stamp.

Every failure prints exactly one line on stderr, 'namo: <message>', and
no usage text; a line break inside the message is shown as \n. All
names in a batch are generated before the first one is written, so a
failed run leaves stdout empty. If writing to stdout fails part-way
through a batch, the names already written stay and namo exits 1.`)
}
