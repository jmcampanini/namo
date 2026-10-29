package cli

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestExecuteExitStatus(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantStatus int
		wantError  string
	}{
		{name: "success", args: []string{"--no-stamp"}},
		{name: "version", args: []string{"--version"}},
		{name: "version ignores flag values", args: []string{"--version", "--count=0"}},
		{name: "root help", args: []string{"--help"}},
		{name: "help flag ignores flag values", args: []string{"--help", "--count=0"}},
		{name: "help flag before topic", args: []string{"--help", "pattern"}},
		{name: "help topic", args: []string{"exit-codes"}},
		{name: "help command", args: []string{"help"}},
		{name: "help command with topic", args: []string{"help", "pattern"}},
		{name: "completion", args: []string{"completion", "zsh"}},
		{name: "help flag on completion", args: []string{"completion", "--help"}},
		{name: "help flag with valid operand", args: []string{"completion", "--help", "fish"}},
		{name: "unknown command", args: []string{"unknown"}, wantStatus: 2, wantError: `unknown command "unknown" for "namo"`},
		{name: "root operand after --", args: []string{"--", "extra"}, wantStatus: 2, wantError: `unknown command "extra" for "namo"`},
		{name: "unknown flag", args: []string{"--unknown"}, wantStatus: 2, wantError: "unknown flag: --unknown"},
		{name: "topic operand", args: []string{"pattern", "extra"}, wantStatus: 2, wantError: `unknown command "extra" for "namo pattern"`},
		{name: "unknown help topic", args: []string{"help", "unknown"}, wantStatus: 2, wantError: `unknown command "unknown" for "namo"`},
		{name: "help operand after topic", args: []string{"help", "pattern", "extra"}, wantStatus: 2, wantError: `unknown command "extra" for "namo pattern"`},
		{name: "missing completion shell", args: []string{"completion"}, wantStatus: 2, wantError: "accepts 1 arg(s), received 0"},
		{name: "unknown completion shell", args: []string{"completion", "unknown"}, wantStatus: 2, wantError: `invalid argument "unknown" for "namo completion"`},
		{name: "help flag with topic operand", args: []string{"pattern", "--help", "extra"}, wantStatus: 2, wantError: `unknown command "extra" for "namo pattern"`},
		{name: "help flag with root operand", args: []string{"--help", "--", "extra"}, wantStatus: 2, wantError: `unknown command "extra" for "namo"`},
		{name: "help flag with unknown help topic", args: []string{"help", "--help", "unknown"}, wantStatus: 2, wantError: `unknown command "unknown" for "namo"`},
		{name: "help flag with unknown completion shell", args: []string{"completion", "--help", "unknown"}, wantStatus: 2, wantError: `invalid argument "unknown" for "namo completion"`},
		{name: "version with root operand", args: []string{"--version", "--", "extra"}, wantStatus: 2, wantError: `unknown command "extra" for "namo"`},
		{name: "unknown flag next to help flag", args: []string{"--help", "--bogus"}, wantStatus: 2, wantError: "unknown flag: --bogus"},
		{name: "help flag with exclusive prefixes", args: []string{"--help", "--prefix=a", "--raw-prefix=b"}, wantStatus: 2, wantError: "none of the others can be"},
		{name: "version with exclusive stamps", args: []string{"--version", "--no-stamp", "--short-stamp"}, wantStatus: 2, wantError: "none of the others can be"},
		{name: "invalid integer", args: []string{"--count", "x"}, wantStatus: 2, wantError: `invalid argument "x"`},
		{name: "invalid boolean", args: []string{"--no-stamp=x"}, wantStatus: 2, wantError: `invalid argument "x"`},
		{name: "missing flag value", args: []string{"--count"}, wantStatus: 2, wantError: "flag needs an argument"},
		{name: "exclusive prefixes", args: []string{"--prefix=a", "--raw-prefix=b"}, wantStatus: 2, wantError: "none of the others can be"},
		{name: "exclusive stamps", args: []string{"--no-stamp", "--short-stamp"}, wantStatus: 2, wantError: "none of the others can be"},
		{name: "invalid count", args: []string{"--count=0"}, wantStatus: 1, wantError: "count must be between 1 and 100"},
		{name: "invalid prefix", args: []string{"--prefix=!!!"}, wantStatus: 1, wantError: "prefix must contain at least one ASCII letter or digit"},
		{name: "invalid size", args: []string{"--size=bogus"}, wantStatus: 1, wantError: `invalid size "bogus"`},
		{name: "invalid stamp verb", args: []string{"--stamp=%q"}, wantStatus: 1, wantError: "unsupported stamp verb %q"},
		{name: "bare percent", args: []string{"--stamp=%"}, wantStatus: 1, wantError: "ends with a bare %"},
		{name: "newline in message", args: []string{"--stamp=%\n"}, wantStatus: 1, wantError: `unsupported stamp verb %\n`},
		{name: "carriage return in message", args: []string{"--stamp=%\r"}, wantStatus: 1, wantError: `unsupported stamp verb %\r`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, stdout, got := runCLI(t, tt.args...)

			if status != tt.wantStatus {
				t.Errorf("execute(%v) = %d, want %d; stderr = %q", tt.args, status, tt.wantStatus, got)
			}
			if tt.wantStatus == 0 {
				if stdout == "" || got != "" {
					t.Errorf("execute(%v) stdout = %q, stderr = %q, want output on stdout only", tt.args, stdout, got)
				}
				return
			}
			if stdout != "" {
				t.Errorf("stdout = %q, want empty", stdout)
			}
			if !strings.Contains(got, tt.wantError) {
				t.Errorf("stderr = %q, want containing %q", got, tt.wantError)
			}
			if !strings.HasPrefix(got, "namo: ") || strings.Count(got, "\n") != 1 || !strings.HasSuffix(got, "\n") || strings.Contains(got, "\r") {
				t.Errorf("stderr = %q, want one 'namo: <message>' line", got)
			}
			for _, unwanted := range []string{"error:", "Error:", "Usage:"} {
				if strings.Contains(got, unwanted) {
					t.Errorf("stderr = %q, contains %q", got, unwanted)
				}
			}
		})
	}
}

func TestExecuteWriteFailure(t *testing.T) {
	for _, args := range [][]string{
		{"--no-stamp"},
		{"--version"},
		{"--help"},
		{"pattern"},
		{"pattern", "--help"},
		{"help", "pattern"},
		{"completion", "bash"},
	} {
		root := newRootCmd()
		var stderr bytes.Buffer
		root.SetOut(closedWriter{})
		root.SetErr(&stderr)
		root.SetArgs(args)

		if status := execute(root); status != 1 {
			t.Errorf("execute(%v) = %d, want 1", args, status)
		}
		if got, want := stderr.String(), "namo: io: read/write on closed pipe\n"; got != want {
			t.Errorf("execute(%v) stderr = %q, want %q", args, got, want)
		}
	}
}

type closedWriter struct{}

func (closedWriter) Write([]byte) (int, error) {
	return 0, io.ErrClosedPipe
}
