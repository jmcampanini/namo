package cli

import (
	"strings"
	"testing"
)

func TestExitCodesTopicListsEveryStatus(t *testing.T) {
	_, stdout, stderr, err := executeCommand(t, newRootCmd(), "help", "exit-codes")
	if err != nil {
		t.Fatalf("Execute(help exit-codes) error = %v, stderr = %s", err, stderr)
	}
	for _, want := range []string{"\n  0  ", "\n  1  ", "\n  2  ", "'namo: <message>'"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("exit-codes help missing %q:\n%s", want, stdout)
		}
	}
}
