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
		{name: "help topic", args: []string{"exit-codes"}},
		{name: "unknown help topic", args: []string{"help", "unknown"}},
		{name: "missing completion shell", args: []string{"completion"}},
		{name: "unknown completion shell", args: []string{"completion", "unknown"}},
		{name: "unknown command", args: []string{"unknown"}, wantStatus: 2, wantError: `unknown command "unknown"`},
		{name: "unknown flag", args: []string{"--unknown"}, wantStatus: 2, wantError: "unknown flag: --unknown"},
		{name: "topic operand", args: []string{"pattern", "extra"}, wantStatus: 2, wantError: `unknown command "extra"`},
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
			root := newRootCmd()
			var stdout, stderr bytes.Buffer
			root.SetOut(&stdout)
			root.SetErr(&stderr)
			root.SetArgs(tt.args)

			status := execute(root)

			if status != tt.wantStatus {
				t.Errorf("execute(%v) = %d, want %d; stderr = %q", tt.args, status, tt.wantStatus, stderr.String())
			}
			if tt.wantStatus == 0 {
				if stdout.Len() == 0 {
					t.Error("stdout is empty, want names or help")
				}
				return
			}
			if stdout.Len() != 0 {
				t.Errorf("stdout = %q, want empty", stdout.String())
			}
			got := stderr.String()
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
	root := newRootCmd()
	var stderr bytes.Buffer
	root.SetOut(closedWriter{})
	root.SetErr(&stderr)
	root.SetArgs([]string{"--no-stamp"})

	if status := execute(root); status != 1 {
		t.Errorf("execute() = %d, want 1", status)
	}
	if got, want := stderr.String(), "namo: io: read/write on closed pipe\n"; got != want {
		t.Errorf("stderr = %q, want %q", got, want)
	}
}

type closedWriter struct{}

func (closedWriter) Write([]byte) (int, error) {
	return 0, io.ErrClosedPipe
}
