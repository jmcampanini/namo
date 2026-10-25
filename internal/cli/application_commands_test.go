package cli

import "github.com/spf13/cobra"

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
