package cli

import (
	"strings"
	"testing"
)

func TestHelpCommandPrintsTheNamedCommandsHelp(t *testing.T) {
	for path := range applicationCommands(newRootCmd()) {
		t.Run(path, func(t *testing.T) {
			args := strings.Fields(path)[1:]
			status, want, stderr := runCLI(t, append(args, "--help")...)
			if status != 0 || stderr != "" {
				t.Fatalf("execute(%v --help) = %d, stderr = %q", args, status, stderr)
			}
			invocation := append([]string{"help"}, args...)
			status, got, stderr := runCLI(t, invocation...)
			if status != 0 || stderr != "" {
				t.Fatalf("execute(%v) = %d, stderr = %q", invocation, status, stderr)
			}
			if got != want {
				t.Errorf("execute(%v) output differs from --help:\n%s\n---\n%s", invocation, got, want)
			}
		})
	}
}

func TestHelpCommandCompletesCommandNames(t *testing.T) {
	_, stdout, _, err := executeCommand(t, newRootCmd(), "__complete", "help", "")
	if err != nil {
		t.Fatalf("Execute(__complete help) error = %v", err)
	}
	for _, name := range []string{"completion", "exit-codes", "help", "pattern", "stamps"} {
		if !strings.Contains(stdout, "\n"+name+"\t") && !strings.HasPrefix(stdout, name+"\t") {
			t.Errorf("help completions do not offer %q:\n%s", name, stdout)
		}
	}
	if strings.Contains(stdout, "__complete") {
		t.Errorf("help completions offer a hidden command:\n%s", stdout)
	}
}
