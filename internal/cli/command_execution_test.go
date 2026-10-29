package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(t *testing.T, root *cobra.Command, args ...string) (*cobra.Command, string, string, error) {
	t.Helper()
	if args == nil {
		args = []string{}
	}
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SetArgs(args)
	command, err := root.ExecuteC()
	return command, stdout.String(), stderr.String(), err
}

// runCLI runs args through execute, the path the binary takes, and returns
// the exit status, stdout, and stderr.
func runCLI(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	if args == nil {
		args = []string{}
	}
	root := newRootCmd()
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SetArgs(args)
	status := execute(root)
	return status, stdout.String(), stderr.String()
}
