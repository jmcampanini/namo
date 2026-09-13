package cmd

import (
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
