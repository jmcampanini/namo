package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

type grammarCoverage uint8

const (
	grammarCoverageNone grammarCoverage = iota
	grammarCoverageValid
	grammarCoverageRejectedOperand
)

type grammarCase struct {
	name       string
	args       []string
	path       string
	coverage   grammarCoverage
	wantRunner string
	wantErr    string
	wantOutput bool
}

var grammarCases = []grammarCase{
	{
		name:       "bare root",
		path:       "namo",
		coverage:   grammarCoverageValid,
		wantRunner: "namo",
		wantOutput: true,
	},
	{
		name:     "root rejects operand",
		args:     []string{"extra"},
		path:     "namo",
		coverage: grammarCoverageRejectedOperand,
		wantErr:  "unknown command",
	},
	{
		name:    "root rejects unknown flag",
		args:    []string{"--bogus"},
		path:    "namo",
		wantErr: "unknown flag",
	},
	{
		name:       "root help",
		args:       []string{"--help"},
		path:       "namo",
		wantOutput: true,
	},
	{
		name:       "root version",
		args:       []string{"--version"},
		path:       "namo",
		wantOutput: true,
	},
	{
		name:       "docs",
		args:       []string{"docs"},
		path:       "namo docs",
		coverage:   grammarCoverageValid,
		wantRunner: "namo docs",
		wantOutput: true,
	},
	{
		name:     "docs rejects operand",
		args:     []string{"docs", "extra"},
		path:     "namo docs",
		coverage: grammarCoverageRejectedOperand,
		wantErr:  "unknown command",
	},
	{
		name:    "docs rejects unknown flag",
		args:    []string{"docs", "--bogus"},
		path:    "namo docs",
		wantErr: "unknown flag",
	},
	{
		name:       "docs help",
		args:       []string{"docs", "--help"},
		path:       "namo docs",
		wantOutput: true,
	},
	{
		name:       "completion bash",
		args:       []string{"completion", "bash"},
		path:       "namo completion bash",
		wantOutput: true,
	},
}

func TestApplicationCommandGrammarInventory(t *testing.T) {
	allCommands, applicationCommands := commandInventories()
	coverageByPath := make(map[string]struct {
		valid           bool
		rejectedOperand bool
	})

	for _, test := range grammarCases {
		if _, ok := allCommands[test.path]; !ok {
			t.Errorf("grammar case %q targets %q, which is not in the command tree", test.name, test.path)
		}
		if test.coverage == grammarCoverageNone {
			continue
		}
		if _, ok := applicationCommands[test.path]; !ok {
			t.Errorf("grammar case %q provides application coverage for %q, which is not an application command", test.name, test.path)
			continue
		}

		coverage := coverageByPath[test.path]
		switch test.coverage {
		case grammarCoverageValid:
			if test.wantErr != "" || test.wantRunner != test.path {
				t.Errorf("valid grammar case %q must succeed and invoke %q", test.name, test.path)
				continue
			}
			coverage.valid = true
		case grammarCoverageRejectedOperand:
			if test.wantErr == "" || test.wantRunner != "" {
				t.Errorf("rejected-operand grammar case %q must fail before invoking a runner", test.name)
				continue
			}
			coverage.rejectedOperand = true
		default:
			t.Errorf("grammar case %q has unknown coverage classification %d", test.name, test.coverage)
			continue
		}
		coverageByPath[test.path] = coverage
	}

	for path, command := range applicationCommands {
		if command.Args == nil {
			t.Errorf("%s has no explicit Args validator", path)
		}
		coverage := coverageByPath[path]
		if !coverage.valid {
			t.Errorf("%s has no valid grammar case", path)
		}
		if !coverage.rejectedOperand {
			t.Errorf("%s has no rejected-operand grammar case", path)
		}
	}
}

func TestCommandGrammar(t *testing.T) {
	for _, test := range grammarCases {
		t.Run(test.name, func(t *testing.T) {
			root := newRootCmd()
			invocations := spyCommandRunners(root)
			command, stdout, stderr, err := executeCommand(t, root, test.args...)

			if command == nil {
				t.Fatalf("Execute(%v) selected no command", test.args)
			}
			if path := command.CommandPath(); path != test.path {
				t.Errorf("Execute(%v) selected %q, want %q", test.args, path, test.path)
			}
			if test.wantErr == "" {
				if err != nil {
					t.Fatalf("Execute(%v) error = %v", test.args, err)
				}
			} else if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("Execute(%v) error = %v, want containing %q", test.args, err, test.wantErr)
			}
			if test.wantErr != "" && stderr != "" {
				t.Errorf("Execute(%v) stderr = %q, want empty", test.args, stderr)
			}

			if test.wantOutput && stdout == "" {
				t.Errorf("Execute(%v) stdout is empty, want output", test.args)
			}
			if !test.wantOutput && stdout != "" {
				t.Errorf("Execute(%v) stdout = %q, want empty", test.args, stdout)
			}
			assertRunnerInvocations(t, invocations, test.wantRunner)
		})
	}
}

func commandInventories() (map[string]*cobra.Command, map[string]*cobra.Command) {
	root := newRootCmd()
	applicationCommands := make(map[string]*cobra.Command)
	collectCommands(root, applicationCommands)

	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()
	allCommands := make(map[string]*cobra.Command)
	collectCommands(root, allCommands)
	return allCommands, applicationCommands
}

func collectCommands(command *cobra.Command, commands map[string]*cobra.Command) {
	commands[command.CommandPath()] = command
	for _, child := range command.Commands() {
		collectCommands(child, commands)
	}
}

func spyCommandRunners(command *cobra.Command) map[string]int {
	invocations := make(map[string]int)
	var walk func(*cobra.Command)
	walk = func(current *cobra.Command) {
		if current.RunE != nil {
			path := current.CommandPath()
			runE := current.RunE
			invocations[path] = 0
			current.RunE = func(command *cobra.Command, args []string) error {
				invocations[path]++
				return runE(command, args)
			}
		}
		for _, child := range current.Commands() {
			walk(child)
		}
	}
	walk(command)
	return invocations
}

func assertRunnerInvocations(t *testing.T, invocations map[string]int, wantPath string) {
	t.Helper()
	if wantPath != "" {
		if _, ok := invocations[wantPath]; !ok {
			t.Fatalf("runner %q was not spied", wantPath)
		}
	}
	for path, got := range invocations {
		want := 0
		if path == wantPath {
			want = 1
		}
		if got != want {
			t.Errorf("runner %q invocations = %d, want %d", path, got, want)
		}
	}
}
