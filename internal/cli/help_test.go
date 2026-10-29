package cli

import (
	"strings"
	"testing"
)

func TestHelpCommandPrintsTheNamedCommandsHelp(t *testing.T) {
	for path := range applicationCommands(newRootCmd()) {
		t.Run(path, func(t *testing.T) {
			args := strings.Fields(path)[1:]
			_, want, _, err := executeCommand(t, newRootCmd(), append(args, "--help")...)
			if err != nil {
				t.Fatalf("Execute(%v --help) error = %v", args, err)
			}
			invocation := append([]string{"help"}, args...)
			_, got, stderr, err := executeCommand(t, newRootCmd(), invocation...)
			if err != nil || stderr != "" {
				t.Fatalf("Execute(%v) error = %v, stderr = %q", invocation, err, stderr)
			}
			if got != want {
				t.Errorf("Execute(%v) output differs from --help:\n%s\n---\n%s", invocation, got, want)
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
