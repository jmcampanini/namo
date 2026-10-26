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
// Cobra owns its grammar and it declares no validator.
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
		if !command.HasSubCommands() && !isHelpTopic(command) {
			continue
		}
		if !command.Runnable() {
			t.Errorf("%s is a command group or help topic but is not runnable", path)
		}
		if command.Args == nil || reflect.ValueOf(command.Args).Pointer() != noArgs {
			t.Errorf("%s is a command group or help topic but its Args is not cobra.NoArgs", path)
		}
	}
}
