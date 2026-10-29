package cli

import "github.com/spf13/cobra"

// applicationCommands maps command paths to every command in root's tree.
func applicationCommands(root *cobra.Command) map[string]*cobra.Command {
	commands := map[string]*cobra.Command{root.CommandPath(): root}
	var walk func(*cobra.Command)
	walk = func(command *cobra.Command) {
		for _, child := range command.Commands() {
			commands[child.CommandPath()] = child
			walk(child)
		}
	}
	walk(root)
	return commands
}

// isHelpTopic reports whether command is a help topic, which newHelpTopic
// places in the help topic group.
func isHelpTopic(command *cobra.Command) bool {
	return command.GroupID == helpTopicGroupID
}
