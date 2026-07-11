package cmd

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func runCommand(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	if args == nil {
		args = []string{}
	}
	root := newRootCmd()
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SetArgs(args)
	err := root.Execute()
	return stdout.String(), stderr.String(), err
}

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantLines int
		wantMatch string
		wantErr   string
	}{
		{name: "default", args: []string{}, wantLines: 1, wantMatch: `^\d{12}-[a-z0-9]+(-[a-z0-9]+)+$`},
		{name: "prefix", args: []string{"-p", "debug-output"}, wantLines: 1, wantMatch: `^debug-output-\d{12}-[a-z0-9-]+$`},
		{name: "no stamp", args: []string{"--no-stamp"}, wantLines: 1, wantMatch: `^[a-z0-9]+(-[a-z0-9]+)+$`},
		{name: "short stamp", args: []string{"--short-stamp"}, wantLines: 1, wantMatch: `^\d{4}-[a-z0-9-]+$`},
		{name: "custom stamp", args: []string{"--stamp", "%Y%m%d"}, wantLines: 1, wantMatch: `^\d{8}-[a-z0-9-]+$`},
		{name: "count", args: []string{"-n", "3"}, wantLines: 3, wantMatch: `^\d{12}-[a-z0-9-]+$`},
		{name: "size short", args: []string{"-s", "short", "--no-stamp"}, wantLines: 1, wantMatch: `^[a-z0-9]+(-[a-z0-9]+)+$`},
		{name: "unsupported stamp verb", args: []string{"--stamp", "%q"}, wantErr: "unsupported stamp verb"},
		{name: "stamp flags mutually exclusive", args: []string{"--no-stamp", "--short-stamp"}, wantErr: "none of the others can be"},
		{name: "custom stamp with no-stamp rejected", args: []string{"--stamp", "%y", "--no-stamp"}, wantErr: "none of the others can be"},
		{name: "invalid size", args: []string{"-s", "bogus"}, wantErr: "invalid size"},
		{name: "zero count", args: []string{"-n", "0"}, wantErr: "count must be at least 1"},
		{name: "positional arg rejected", args: []string{"extra"}, wantErr: "unknown command"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, _, err := runCommand(t, tt.args...)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Execute(%v) error = %v, want containing %q", tt.args, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Execute(%v) error = %v", tt.args, err)
			}
			lines := strings.Split(strings.TrimSuffix(stdout, "\n"), "\n")
			if len(lines) != tt.wantLines {
				t.Fatalf("Execute(%v) printed %d lines, want %d:\n%s", tt.args, len(lines), tt.wantLines, stdout)
			}
			match := regexp.MustCompile(tt.wantMatch)
			for _, line := range lines {
				if !match.MatchString(line) {
					t.Fatalf("line %q does not match %q", line, tt.wantMatch)
				}
			}
		})
	}
}

func TestRootCommandBatchSharesStamp(t *testing.T) {
	stdout, _, err := runCommand(t, "-n", "5")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	lines := strings.Split(strings.TrimSuffix(stdout, "\n"), "\n")
	if len(lines) != 5 {
		t.Fatalf("printed %d lines, want 5", len(lines))
	}
	stamps := make(map[string]struct{})
	names := make(map[string]struct{})
	for _, line := range lines {
		stamps[strings.SplitN(line, "-", 2)[0]] = struct{}{}
		names[line] = struct{}{}
	}
	if len(stamps) != 1 {
		t.Fatalf("batch used %d distinct stamps, want 1:\n%s", len(stamps), stdout)
	}
	if len(names) != 5 {
		t.Fatalf("batch has %d distinct names, want 5:\n%s", len(names), stdout)
	}
}

func TestDocsCommand(t *testing.T) {
	stdout, _, err := runCommand(t, "docs")
	if err != nil {
		t.Fatalf("Execute(docs) error = %v", err)
	}
	if !strings.Contains(stdout, "namo") {
		t.Fatalf("docs output does not mention namo:\n%s", stdout)
	}
}

func TestVersionFlag(t *testing.T) {
	stdout, _, err := runCommand(t, "--version")
	if err != nil {
		t.Fatalf("Execute(--version) error = %v", err)
	}
	if !strings.Contains(stdout, "namo version") {
		t.Fatalf("version output = %q, want containing %q", stdout, "namo version")
	}
}
