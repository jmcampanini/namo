package cli

import (
	"strings"
	"testing"
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

func firstLine(s string) string {
	line, _, _ := strings.Cut(s, "\n")
	return line
}
