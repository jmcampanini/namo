package cli

import "github.com/spf13/cobra"

// topicContract is the output contract every help topic shares, composed
// into each topic's long description so the topics cannot drift apart.
const topicContract = `This topic prints plain text to stdout and nothing to stderr, never
prompts, works without a terminal, and reads and writes no files.`

// newHelpTopic builds a help topic: a command that does no work and prints
// the same text whether invoked directly or through 'namo help <name>'.
func newHelpTopic(use, short, long string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long + "\n\n" + topicContract,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
}
