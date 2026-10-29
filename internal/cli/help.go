package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newHelpCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "help [command]",
		Short:   "Print the help of a command or help topic",
		GroupID: commandGroupID,
		Long: `'namo help' prints the root help, the same text as 'namo --help', and
'namo help NAME' prints the help of the command or help topic NAME, the
same text as 'namo NAME --help'. A NAME that is not a command or help
topic, or an operand after one, is a usage error.

Help goes to stdout as plain text, the only output format. 'namo help'
never prompts, works without a terminal, and reads and writes no files.

` + failureContract,
		Args: func(cmd *cobra.Command, args []string) error {
			_, err := findHelpTarget(cmd, args)
			return err
		},
		ValidArgsFunction: completeHelpTarget,
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := findHelpTarget(cmd, args)
			if err != nil {
				return err
			}
			// Cobra adds the help flag when a command runs; add it here so
			// the help matches what 'namo NAME --help' prints.
			target.InitDefaultHelpFlag()
			return target.Help()
		},
	}
}

// findHelpTarget returns the command or help topic that args name, walking
// only the command tree.
func findHelpTarget(cmd *cobra.Command, args []string) (*cobra.Command, error) {
	target, rest, err := cmd.Root().Find(args)
	if err != nil {
		return nil, err
	}
	if len(rest) > 0 {
		return nil, fmt.Errorf("unknown command %q for %q", rest[0], target.CommandPath())
	}
	return target, nil
}

// completeHelpTarget completes the names of the commands and help topics
// under the command that args already name.
func completeHelpTarget(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	target, err := findHelpTarget(cmd, args)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var completions []cobra.Completion
	for _, command := range target.Commands() {
		if command.Hidden || !strings.HasPrefix(command.Name(), toComplete) {
			continue
		}
		completions = append(completions, cobra.CompletionWithDesc(command.Name(), command.Short))
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}
