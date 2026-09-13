package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestEveryApplicationCommandDeclaresPositionalGrammar requires an Args
// validator on every application-owned command and a runner on every command
// with subcommands. Cobra validates operands before any runner, so this is
// enough evidence that rejected operands reach no command work.
func TestEveryApplicationCommandDeclaresPositionalGrammar(t *testing.T) {
	root := newRootCmd()

	// The root is a leaf with its own grammar, so it declares an arity too.
	if root.Args == nil {
		t.Errorf("%s has no Args validator", root.CommandPath())
	}
	for path, command := range applicationCommands(root) {
		if command != root && command.Args == nil {
			t.Errorf("%s has no Args validator", path)
		}
		if command.HasSubCommands() && command.RunE == nil {
			t.Errorf("%s has subcommands but no RunE", path)
		}
	}
}

// TestUnknownFlagIsRejectedWithoutUsageOutput checks that SilenceUsage and
// SilenceErrors on the root reach every command: an unknown flag returns an
// error and writes nothing, leaving main to print the single error line.
func TestUnknownFlagIsRejectedWithoutUsageOutput(t *testing.T) {
	for _, args := range [][]string{{"--bogus"}, {"docs", "--bogus"}, {"exit-codes", "--bogus"}} {
		_, stdout, stderr, err := executeCommand(t, newRootCmd(), args...)

		if err == nil || !strings.Contains(err.Error(), "unknown flag") {
			t.Errorf("Execute(%v) error = %v, want containing %q", args, err, "unknown flag")
		}
		if stdout != "" || stderr != "" {
			t.Errorf("Execute(%v) stdout = %q, stderr = %q, want both empty", args, stdout, stderr)
		}
	}
}

// applicationCommands maps command paths to every command in root's tree
// except Cobra's built-in help and completion commands.
func applicationCommands(root *cobra.Command) map[string]*cobra.Command {
	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()

	commands := map[string]*cobra.Command{root.CommandPath(): root}
	var walk func(*cobra.Command)
	walk = func(command *cobra.Command) {
		for _, child := range command.Commands() {
			if child.Name() == "help" || child.Name() == "completion" {
				continue
			}
			commands[child.CommandPath()] = child
			walk(child)
		}
	}
	walk(root)
	return commands
}
