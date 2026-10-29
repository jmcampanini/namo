package cli

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/spf13/cobra"
)

const slugPattern = `[a-z0-9]+(-[a-z0-9]+)*`

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantLines int
		wantMatch string
		wantErr   string
	}{
		{name: "default", args: []string{}, wantLines: 1, wantMatch: `^\d{12}-` + slugPattern + `$`},
		{name: "prefix", args: []string{"-p", "debug-output"}, wantLines: 1, wantMatch: `^debug-output-\d{12}-` + slugPattern + `$`},
		{name: "no stamp", args: []string{"--no-stamp"}, wantLines: 1, wantMatch: `^` + slugPattern + `$`},
		{name: "short stamp", args: []string{"--short-stamp"}, wantLines: 1, wantMatch: `^\d{4}-` + slugPattern + `$`},
		{name: "custom stamp", args: []string{"--stamp", "%Y%m%d"}, wantLines: 1, wantMatch: `^\d{8}-` + slugPattern + `$`},
		{name: "count", args: []string{"-n", "3"}, wantLines: 3, wantMatch: `^\d{12}-` + slugPattern + `$`},
		{name: "size short", args: []string{"-s", "short", "--no-stamp"}, wantLines: 1, wantMatch: `^` + slugPattern + `$`},
		{name: "unsupported stamp verb", args: []string{"--stamp", "%q"}, wantErr: "unsupported stamp verb"},
		{name: "stamp flags mutually exclusive", args: []string{"--no-stamp", "--short-stamp"}, wantErr: "none of the others can be"},
		{name: "custom stamp with no-stamp rejected", args: []string{"--stamp", "%y", "--no-stamp"}, wantErr: "none of the others can be"},
		{name: "invalid size", args: []string{"-s", "bogus"}, wantErr: "invalid size"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, stdout, _, err := executeCommand(t, newRootCmd(), tt.args...)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Execute(%v) error = %v, want containing %q", tt.args, err, tt.wantErr)
				}
				if stdout != "" {
					t.Fatalf("Execute(%v) stdout = %q, want empty", tt.args, stdout)
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

func TestRootCommandCountValidation(t *testing.T) {
	for _, count := range []int{1, 100} {
		_, stdout, stderr, err := executeCommand(t, newRootCmd(), "--count="+strconv.Itoa(count), "--no-stamp")
		if err != nil {
			t.Fatalf("Execute(--count=%d) error = %v", count, err)
		}
		if stderr != "" {
			t.Fatalf("Execute(--count=%d) stderr = %q, want empty", count, stderr)
		}
		if lines := strings.Count(stdout, "\n"); lines != count {
			t.Fatalf("Execute(--count=%d) printed %d lines, want %d", count, lines, count)
		}
	}

	invalid := []struct {
		name string
		args []string
	}{
		{name: "zero", args: []string{"--count=0"}},
		{name: "negative", args: []string{"--count=-1"}},
		{name: "above maximum", args: []string{"--count=101"}},
		{name: "very large", args: []string{"--count=" + strconv.Itoa(int(^uint(0)>>1))}},
		{name: "validated before size", args: []string{"--count=101", "--size=bogus"}},
	}
	for _, tt := range invalid {
		t.Run(tt.name, func(t *testing.T) {
			_, stdout, _, err := executeCommand(t, newRootCmd(), tt.args...)
			if err == nil || !strings.Contains(err.Error(), "count must be between 1 and 100") {
				t.Fatalf("Execute(%v) error = %v, want count range error", tt.args, err)
			}
			if stdout != "" {
				t.Fatalf("Execute(%v) stdout = %q, want empty", tt.args, stdout)
			}
		})
	}
}

func TestRootCommandStrictPrefix(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "complex normalization", args: []string{"--prefix", "  Debug__Output---V2  "}, want: "debug-output-v2"},
		{name: "repeated flag uses last value", args: []string{"-p", "first", "--prefix", "FINAL value"}, want: "final-value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args = append(tt.args, "--no-stamp")
			_, stdout, stderr, err := executeCommand(t, newRootCmd(), tt.args...)
			if err != nil {
				t.Fatalf("Execute(%v) error = %v", tt.args, err)
			}
			if stderr != "" {
				t.Fatalf("Execute(%v) stderr = %q, want empty", tt.args, stderr)
			}
			assertPrefixedName(t, stdout, tt.want)
		})
	}
}

func TestRootCommandStrictPrefixRejectsEmptyNormalizedInput(t *testing.T) {
	for _, prefix := range []string{"", "--- ! ☃ ---"} {
		_, stdout, _, err := executeCommand(t, newRootCmd(), "--prefix", prefix, "--no-stamp")
		if err == nil || !strings.Contains(err.Error(), "prefix must contain at least one ASCII letter or digit") {
			t.Fatalf("Execute(--prefix %q) error = %v, want ASCII validation error", prefix, err)
		}
		if stdout != "" {
			t.Fatalf("Execute(--prefix %q) stdout = %q, want empty", prefix, stdout)
		}
	}
}

func TestRootCommandRawPrefix(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "preserved verbatim", args: []string{"--raw-prefix=--Raw Prefix!?--"}, want: "--Raw Prefix!?--"},
		{name: "unicode preserved", args: []string{"--raw-prefix", "東京💥é"}, want: "東京💥é"},
		{name: "trailing hyphen keeps joining hyphen", args: []string{"--raw-prefix", "trusted-"}, want: "trusted-"},
		{name: "embedded newline preserved", args: []string{"--raw-prefix", "first\nsecond"}, want: "first\nsecond"},
		{name: "repeated flag uses last value", args: []string{"--raw-prefix", "first", "--raw-prefix", "LAST value!?"}, want: "LAST value!?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args = append(tt.args, "--no-stamp")
			_, stdout, stderr, err := executeCommand(t, newRootCmd(), tt.args...)
			if err != nil {
				t.Fatalf("Execute(%v) error = %v", tt.args, err)
			}
			if stderr != "" {
				t.Fatalf("Execute(%v) stderr = %q, want empty", tt.args, stderr)
			}
			assertPrefixedName(t, stdout, tt.want)
		})
	}
}

func TestRootCommandEmptyRawPrefixBehavesAsNoPrefix(t *testing.T) {
	_, stdout, stderr, err := executeCommand(t, newRootCmd(), "--raw-prefix=", "--no-stamp")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if stderr != "" {
		t.Fatalf("Execute() stderr = %q, want empty", stderr)
	}
	if !regexp.MustCompile(`^` + slugPattern + `\n$`).MatchString(stdout) {
		t.Fatalf("output = %q, want an unprefixed slug", stdout)
	}
}

func TestRootCommandCustomStampCanOptOutOfLineSafety(t *testing.T) {
	_, stdout, stderr, err := executeCommand(t, newRootCmd(), "--stamp=trusted\nstamp")
	if err != nil {
		t.Fatalf("Execute(custom newline stamp) error = %v", err)
	}
	if stderr != "" {
		t.Fatalf("Execute(custom newline stamp) stderr = %q, want empty", stderr)
	}
	if !strings.HasPrefix(stdout, "trusted\nstamp-") || strings.Count(stdout, "\n") < 2 {
		t.Fatalf("Execute(custom newline stamp) stdout = %q, want preserved newline", stdout)
	}
}

func TestRootCommandPrefixModesAreMutuallyExclusive(t *testing.T) {
	for _, args := range [][]string{
		{"--prefix", "strict", "--raw-prefix", "raw"},
		{"--prefix=", "--raw-prefix="},
	} {
		_, stdout, _, err := executeCommand(t, newRootCmd(), args...)
		if err == nil {
			t.Fatalf("Execute(%v) error = nil, want mutually exclusive flag error", args)
		}
		if stdout != "" {
			t.Fatalf("Execute(%v) stdout = %q, want empty", args, stdout)
		}
	}
}

func assertPrefixedName(t *testing.T, stdout, prefix string) {
	t.Helper()
	wantPrefix := prefix + "-"
	if !strings.HasPrefix(stdout, wantPrefix) {
		t.Fatalf("output = %q, want prefix %q", stdout, wantPrefix)
	}
	slug := strings.TrimPrefix(stdout, wantPrefix)
	if !regexp.MustCompile(`^` + slugPattern + `\n$`).MatchString(slug) {
		t.Fatalf("output slug %q is not a names-only slug", slug)
	}
}

func TestRootCommandBatchSharesStamp(t *testing.T) {
	_, stdout, _, err := executeCommand(t, newRootCmd(), "-n", "5")
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

func TestRootHelpNamesFlagsAndTopics(t *testing.T) {
	_, stdout, stderr, err := executeCommand(t, newRootCmd(), "--help")
	if err != nil {
		t.Fatalf("Execute(--help) error = %v", err)
	}
	if stderr != "" {
		t.Fatalf("Execute(--help) stderr = %q, want empty", stderr)
	}
	for _, want := range []string{"--prefix", "normalized to lowercase ASCII words", "--raw-prefix", "passed through as given", "trusted input", "1 to 100", "strftime-style stamp layout", "--version"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("Execute(--help) output does not contain %q:\n%s", want, stdout)
		}
	}
}

// TestRootHelpListsHelpTopicsApartFromCommands requires generated root help
// to list every help topic under Help topics and every other command under
// Commands, with no command left for an Additional Commands section.
func TestRootHelpListsHelpTopicsApartFromCommands(t *testing.T) {
	root := newRootCmd()
	_, stdout, _, err := executeCommand(t, root, "--help")
	if err != nil {
		t.Fatalf("Execute(--help) error = %v", err)
	}
	if strings.Contains(stdout, "Additional Commands:") {
		t.Errorf("root help has an Additional Commands section:\n%s", stdout)
	}
	sections := map[string]string{}
	for _, title := range []string{"Commands:", "Help topics:"} {
		_, section, ok := strings.Cut(stdout, "\n"+title+"\n")
		if !ok {
			t.Fatalf("root help has no %q section:\n%s", title, stdout)
		}
		section, _, _ = strings.Cut(section, "\n\n")
		sections[title] = "\n" + section
	}
	for _, command := range root.Commands() {
		if command.Hidden {
			continue
		}
		title := "Commands:"
		if isHelpTopic(command) {
			title = "Help topics:"
		}
		if entry := "\n  " + command.Name() + " "; !strings.Contains(sections[title], entry) {
			t.Errorf("root help does not list %q under %q:\n%s", command.Name(), title, stdout)
		}
	}
}

func TestRootVersion(t *testing.T) {
	for _, flag := range []string{"--version", "-v"} {
		_, stdout, stderr, err := executeCommand(t, newRootCmd(), flag)
		if err != nil || stderr != "" {
			t.Fatalf("Execute(%s) error = %v, stderr = %q", flag, err, stderr)
		}
		if want := "namo version " + Version + "\n"; stdout != want {
			t.Errorf("Execute(%s) stdout = %q, want %q", flag, stdout, want)
		}
	}
}

// TestRootDefinesItsOwnHelpCompletionAndVersion requires namo's own help and
// completion commands and --version flag, so running the root adds none of
// the command library's defaults.
func TestRootDefinesItsOwnHelpCompletionAndVersion(t *testing.T) {
	root := newRootCmd()
	if root.Version != "" {
		t.Errorf("root.Version = %q, want empty", root.Version)
	}
	if root.Flags().Lookup("version") == nil {
		t.Error("root has no --version flag")
	}
	constructed := map[string]*cobra.Command{}
	for _, command := range root.Commands() {
		constructed[command.Name()] = command
	}
	for _, name := range []string{"help", "completion"} {
		if constructed[name] == nil {
			t.Errorf("newRootCmd() has no %s command", name)
		}
	}

	if _, _, _, err := executeCommand(t, root, "--no-stamp"); err != nil {
		t.Fatalf("Execute(--no-stamp) error = %v", err)
	}
	commands := root.Commands()
	if len(commands) != len(constructed) {
		t.Errorf("Execute left %d root commands, want %d", len(commands), len(constructed))
	}
	for _, command := range commands {
		if constructed[command.Name()] != command {
			t.Errorf("Execute added or replaced the %s command", command.Name())
		}
	}
}

// TestEveryApplicationCommandMeetsHelpBaseline walks a freshly constructed
// tree and fails on a missing short help, a missing long description, a
// rendered help line over 80 columns or holding a tab, which renders wider
// than it counts, or a help topic whose direct invocation prints text other
// than its help.
func TestEveryApplicationCommandMeetsHelpBaseline(t *testing.T) {
	root := newRootCmd()
	topics := 0
	for path, command := range applicationCommands(root) {
		if isHelpTopic(command) {
			topics++
		}
		t.Run(path, func(t *testing.T) {
			if strings.TrimSpace(command.Short) == "" {
				t.Error("missing short help")
			}
			if strings.TrimSpace(command.Long) == "" {
				t.Error("missing long description")
			}

			args := strings.Fields(path)[1:]
			_, help, stderr, err := executeCommand(t, newRootCmd(), append(args, "--help")...)
			if err != nil || stderr != "" {
				t.Fatalf("Execute(%v --help) error = %v, stderr = %q", args, err, stderr)
			}
			for i, line := range strings.Split(help, "\n") {
				if width := utf8.RuneCountInString(line); width > 80 || strings.Contains(line, "\t") {
					t.Errorf("help line %d is %d columns or holds a tab, want at most 80 columns: %q", i+1, width, line)
				}
			}
			if !isHelpTopic(command) {
				return
			}

			for _, invocation := range [][]string{args, append([]string{"help"}, args...)} {
				_, stdout, stderr, err := executeCommand(t, newRootCmd(), invocation...)
				if err != nil || stderr != "" {
					t.Fatalf("Execute(%v) error = %v, stderr = %q", invocation, err, stderr)
				}
				if stdout != help {
					t.Errorf("Execute(%v) output differs from its help:\n%s\n---\n%s", invocation, stdout, help)
				}
			}
		})
	}
	if topics == 0 {
		t.Error("found no help topic; the exit-codes topic is required")
	}
}
