package cmd

import "github.com/spf13/cobra"

func newExitCodesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "exit-codes",
		Short: "Exit codes and error categories",
		Long: `namo exits 0 or 1:

  0  Success. The requested names are on stdout, one per line, and
     stderr is empty. --help, --version, 'namo docs', 'namo help', and
     this topic also exit 0. 'namo help NAME' with a NAME that is not a
     command prints the root help and exits 0, and 'namo completion'
     with a missing or unknown shell prints the completion help and
     exits 0.
  1  Any failure. namo prints 'namo: error: <message>' on stderr and no
     usage text. Usage errors: an unknown command, operand, or flag, a
     --count that is not an integer, --prefix together with
     --raw-prefix, or more than one of --stamp, --short-stamp, and
     --no-stamp. Input errors: --count outside 1 to 100, a --prefix
     with no ASCII letter or digit, a --size other than short,
     standard, or long, or a --stamp layout with an unsupported verb or
     a trailing bare %. Generation errors: the slug source cannot
     supply --count distinct slugs within 100 rounds.

All names in a batch are generated before the first one is written, so a
failed run leaves stdout empty. If writing to stdout fails part-way
through a batch, the names already written stay and namo exits 1.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
}
