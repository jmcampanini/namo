package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	var noDescriptions bool
	cmd := &cobra.Command{
		Use:     "completion bash|fish|powershell|zsh",
		Short:   "Print a shell completion script",
		GroupID: commandGroupID,
		Long: `'namo completion SHELL' prints the completion script for SHELL, one of
bash, fish, powershell, or zsh, to stdout. The script is shell source
text, the only output format; loading it lets the shell complete namo's
commands, flags, and help topics. The bash script needs the
bash-completion package, and the zsh script needs compinit loaded first.
--no-descriptions leaves the command and flag descriptions out of the
completions. A missing or unknown SHELL is a usage error.

` + textCommandContract,
		Example: `  # Load completions in the current shell
  namo completion fish | source
  namo completion powershell | Out-String | Invoke-Expression

  # Load bash completions, which need bash 4 or later and the
  # bash-completion package loaded first
  source <(namo completion bash)

  # Load zsh completions, which need compinit loaded first
  autoload -U compinit && compinit
  source <(namo completion zsh)

  # Load fish completions in every new session
  namo completion fish > ~/.config/fish/completions/namo.fish`,
		ValidArgs: []string{"bash", "fish", "powershell", "zsh"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			return writeCompletion(cmd.Root(), cmd.OutOrStdout(), args[0], !noDescriptions)
		},
	}
	cmd.Flags().BoolVar(&noDescriptions, "no-descriptions", false, "disable completion descriptions")
	// Offer no file names once the shell operand is given; any further
	// operand is a usage error.
	cmd.CompletionOptions.SetDefaultShellCompDirective(cobra.ShellCompDirectiveNoFileComp)
	return cmd
}

// writeCompletion writes root's completion script for shell to out.
func writeCompletion(root *cobra.Command, out io.Writer, shell string, descriptions bool) error {
	switch shell {
	case "bash":
		return root.GenBashCompletionV2(out, descriptions)
	case "fish":
		return root.GenFishCompletion(out, descriptions)
	case "powershell":
		if descriptions {
			return root.GenPowerShellCompletionWithDesc(out)
		}
		return root.GenPowerShellCompletion(out)
	case "zsh":
		if descriptions {
			return root.GenZshCompletion(out)
		}
		return root.GenZshCompletionNoDesc(out)
	default:
		return fmt.Errorf("unsupported shell %q", shell)
	}
}
