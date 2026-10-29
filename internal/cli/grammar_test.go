package cli

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

// TestEveryApplicationCommandDeclaresPositionalGrammar requires an Args
// validator on every application-owned command other than the root, and a
// runner with cobra.NoArgs on every command group and help topic. Cobra
// validates operands before any runner, so this is enough evidence that
// rejected operands reach no command work. The root has subcommands, so
// Cobra owns its grammar and it declares no validator. A leaf that takes no
// operands and no flags of its own is a help topic, so it must sit in the
// help topic group that the topic checks select.
func TestEveryApplicationCommandDeclaresPositionalGrammar(t *testing.T) {
	root := newRootCmd()
	noArgs := reflect.ValueOf(cobra.NoArgs).Pointer()

	for path, command := range applicationCommands(root) {
		if command == root {
			continue
		}
		if command.Args == nil {
			t.Errorf("%s has no Args validator", path)
		}
		noOperands := reflect.ValueOf(command.Args).Pointer() == noArgs
		if !command.HasSubCommands() && !command.Flags().HasFlags() && noOperands && !isHelpTopic(command) {
			t.Errorf("%s takes no operands or flags of its own but is not in the help topic group", path)
		}
		if !command.HasSubCommands() && !isHelpTopic(command) {
			continue
		}
		if !command.Runnable() {
			t.Errorf("%s is a command group or help topic but is not runnable", path)
		}
		if !noOperands {
			t.Errorf("%s is a command group or help topic but its Args is not cobra.NoArgs", path)
		}
	}
}
