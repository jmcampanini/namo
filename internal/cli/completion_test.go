package cli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCompletionPrintsTheShellScript(t *testing.T) {
	tests := []struct {
		shell  string
		header string
	}{
		{shell: "bash", header: "# bash completion V2 for namo"},
		{shell: "fish", header: "# fish completion for namo"},
		{shell: "powershell", header: "# powershell completion for namo"},
		{shell: "zsh", header: "#compdef namo"},
	}

	for _, tt := range tests {
		for _, noDescriptions := range []bool{false, true} {
			args := []string{"completion", tt.shell}
			if noDescriptions {
				args = append(args, "--no-descriptions")
			}
			t.Run(strings.Join(args, " "), func(t *testing.T) {
				_, stdout, stderr, err := executeCommand(t, newRootCmd(), args...)
				if err != nil || stderr != "" {
					t.Fatalf("Execute(%v) error = %v, stderr = %q", args, err, stderr)
				}
				if !strings.HasPrefix(stdout, tt.header) {
					t.Errorf("Execute(%v) script starts %q, want %q", args, firstLine(stdout), tt.header)
				}
				// The script requests completions without descriptions
				// through Cobra's __completeNoDesc command.
				if got := strings.Contains(stdout, "__completeNoDesc"); got != noDescriptions {
					t.Errorf("Execute(%v) requests completions without descriptions = %t, want %t", args, got, noDescriptions)
				}
			})
		}
	}
}

// TestShellCompletionOffersNoFileNames requires shell completion to offer
// no file names, because no namo command or flag takes a file.
func TestShellCompletionOffersNoFileNames(t *testing.T) {
	want := fmt.Sprintf(":%d\n", cobra.ShellCompDirectiveNoFileComp)
	for _, args := range [][]string{
		{"completion", "bash", ""},
		{"help", "pattern", ""},
		{"pattern", ""},
		{"stamps", ""},
		{"exit-codes", ""},
		{"--no-stamp", ""},
		{"--size", ""},
		{"-n", "3", ""},
		{"--help", ""},
		{"--version", ""},
		{"-v", ""},
	} {
		request := append([]string{"__complete"}, args...)
		_, stdout, _, err := executeCommand(t, newRootCmd(), request...)
		if err != nil {
			t.Fatalf("Execute(%q) error = %v", request, err)
		}
		if stdout != want {
			t.Errorf("Execute(%q) = %q, want %q", request, stdout, want)
		}
	}
}

func firstLine(s string) string {
	line, _, _ := strings.Cut(s, "\n")
	return line
}
