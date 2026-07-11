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
		Long:  "Print the complete namo reference: the name pattern, every flag,\nstamp layout verbs, sizes, and usage recipes. Pipeable to less or bat.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprint(cmd.OutOrStdout(), manual)
			return err
		},
	}
}
