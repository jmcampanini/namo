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
