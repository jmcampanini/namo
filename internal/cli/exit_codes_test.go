package cli

import (
	"strings"
	"testing"
)

func TestExitCodesTopicPrintsSameHelpFromBothEntryPoints(t *testing.T) {
	_, direct, stderr, err := executeCommand(t, newRootCmd(), "exit-codes")
	if err != nil {
		t.Fatalf("Execute(exit-codes) error = %v, stderr = %s", err, stderr)
	}
	_, viaHelp, stderr, err := executeCommand(t, newRootCmd(), "help", "exit-codes")
	if err != nil {
		t.Fatalf("Execute(help exit-codes) error = %v, stderr = %s", err, stderr)
	}

	if direct != viaHelp {
		t.Fatalf("exit-codes output differs between entry points:\n%s\n---\n%s", direct, viaHelp)
	}
	for _, want := range []string{"\n  0  ", "\n  1  ", "namo: error:"} {
		if !strings.Contains(direct, want) {
			t.Fatalf("exit-codes help missing %q:\n%s", want, direct)
		}
	}
}
