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

// isHelpTopic reports whether command is a help topic: a non-root command
// with no subcommands and no flags of its own. Judge a freshly constructed
// tree, because Cobra adds the help flag to a command once it runs.
func isHelpTopic(command *cobra.Command) bool {
	return command.HasParent() && !command.HasSubCommands() && !command.Flags().HasFlags()
}
