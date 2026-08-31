package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed manual.txt
var manual string

func newDocsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "docs",
		Short: "Print the namo reference manual.",
		Long: `Print the namo reference manual: the name pattern, the generation
flags, stamp layout verbs, slug sizes, and usage recipes. The manual is
plain text on stdout with no terminal escapes, so it pipes into less or
bat. Command help ('namo --help' and 'namo help exit-codes') is the
canonical contract; the manual supplements it with longer explanations
and examples.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprint(cmd.OutOrStdout(), manual)
			return err
		},
	}
}
