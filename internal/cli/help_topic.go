package cli

import "github.com/spf13/cobra"

// topicContract is the output contract every help topic shares, composed
// into each topic's long description so the topics cannot drift apart.
const topicContract = `This topic prints plain text to stdout, the only output format.

` + textCommandContract

// helpTopicAnnotation marks each command newHelpTopic builds.
const helpTopicAnnotation = "namo.help-topic"

// newHelpTopic builds a help topic: a command that does no work and prints
// the same text whether invoked directly or through 'namo help <name>'.
// Root help lists it under Help topics, apart from the commands.
func newHelpTopic(use, short, long string) *cobra.Command {
	return &cobra.Command{
		Use:         use,
		Short:       short,
		Long:        long + "\n\n" + topicContract,
		GroupID:     helpTopicGroupID,
		Annotations: map[string]string{helpTopicAnnotation: "true"},
		Args:        cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return writeHelp(cmd, cmd.HelpFunc())
		},
	}
}
