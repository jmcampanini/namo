package cli

import (
	"strings"
	"testing"
)

func TestDocsCommand(t *testing.T) {
	_, stdout, stderr, err := executeCommand(t, newRootCmd(), "docs")
	if err != nil {
		t.Fatalf("Execute(docs) error = %v", err)
	}
	if stdout != manual {
		t.Fatalf("Execute(docs) output does not match embedded manual")
	}
	for _, want := range []string{"[prefix-][stamp-]slug", "Omit the timestamp while retaining any prefix", "non-empty raw prefix", "empty raw prefix behaves as no prefix", "1 <= count <= 100", "custom --stamp is also trusted input"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("Execute(docs) output does not contain %q", want)
		}
	}
	if stderr != "" {
		t.Fatalf("Execute(docs) stderr = %q, want empty", stderr)
	}
}
