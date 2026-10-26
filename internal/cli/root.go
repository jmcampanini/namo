// Package cli wires the cobra commands for the namo binary.
package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jmcampanini/namo/internal/namegen"
)

type rootFlags struct {
	count      int
	noStamp    bool
	prefix     string
	rawPrefix  string
	shortStamp bool
	size       string
	stamp      string
}

// Version is set at build time via -ldflags.
var Version = "n/a"

// Execute parses the command line, runs the requested command, and returns
// the process exit status: 0 on success, 1 when a command's runner fails,
// and 2 on a usage error the command library rejects before any runner. A
// failure prints one line, "namo: <message>", to stderr.
func Execute() int {
	return execute(newRootCmd())
}

// execute runs root with the arguments and streams set on it and returns the
// exit status. A line break inside a message is shown as \n so the
// diagnostic stays one line.
func execute(root *cobra.Command) int {
	err := root.Execute()
	if err == nil {
		return 0
	}
	message := strings.NewReplacer("\r", `\r`, "\n", `\n`).Replace(err.Error())
	_, _ = fmt.Fprintf(root.ErrOrStderr(), "namo: %s\n", message)
	var failure *runError
	if errors.As(err, &failure) {
		return 1
	}
	return 2
}

// runError marks a failure a command's runner returned, so execute can tell
// it from the usage errors the command library returns before any runner.
type runError struct {
	err error
}

func (e *runError) Error() string { return e.err.Error() }

func (e *runError) Unwrap() error { return e.err }

// run wraps a runner so every error it returns is a runError.
func run(runner func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := runner(cmd, args); err != nil {
			return &runError{err: err}
		}
		return nil
	}
}

func newRootCmd() *cobra.Command {
	flags := &rootFlags{}
	root := &cobra.Command{
		Use:   "namo",
		Short: "Generate memorable, sortable names.",
		Long: `Generate memorable, sortable names of the form [prefix-][stamp-]slug, for
example debug-output-260711154501-star-studded-booze-cruise. The stamp is
the local time, yymmddhhmmss by default, so names sort by creation time.
The slug is random words from hotdiva2000: lowercase ASCII letters and
digits joined by single hyphens.

'namo' prints one name. -p/--prefix TEXT adds a prefix: ASCII letters and
digits are kept, letters are lowercased, each run of other characters
becomes one hyphen, edge hyphens are dropped, and a value with no ASCII
letter or digit is an error. --raw-prefix TEXT adds the bytes exactly as
given followed by one hyphen; an empty value adds nothing. -n/--count N
prints N names, 1 to 100, that share one stamp and repeat no slug.
-s/--size short|standard|long sets slug length. --stamp LAYOUT renders a
strftime-style layout, --short-stamp is --stamp %H%M, and --no-stamp
omits the stamp and keeps any prefix. --prefix and --raw-prefix are
mutually exclusive, as are --stamp, --short-stamp, and --no-stamp.

namo has no configuration file, reads no files, writes no files, and
reads nothing from stdin; every option is a flag. The stamp follows the
local time zone.

Names go to stdout, one per line, and that plain text is the only output
format; help and version output go to stdout too. A failure prints one
line on stderr, 'namo: <message>', with no usage text. namo never
prompts, works without a terminal, runs no external program, and never
accesses the network. Without --raw-prefix or a custom --stamp, every
name consists only of lowercase ASCII letters, digits, and hyphens, so
it is safe inside "$(namo ...)" and in file paths. --raw-prefix and
--stamp are trusted input: line breaks, control bytes, and path
separators in them pass through unchanged.

Help topics:
  pattern     The name parts and slug sizes:      namo help pattern
  stamps      Stamp layout verbs and rendering:   namo help stamps
  exit-codes  Exit statuses and the stderr line:  namo help exit-codes`,
		Example: `  # Sortable log or scratch file
  log=".sandbox/$(namo -p fish-profile).log"

  # Family of output files sharing one run prefix
  prefix="$(namo)"
  response="out/${prefix}-response.md"
  meta="out/${prefix}-meta.json"

  # Branch or session handle with no stamp
  git switch -c "spike/$(namo --no-stamp -s short)"

  # Ephemeral name that only needs to be unique today
  namo --short-stamp

  # Batch of seed identifiers
  namo -n 25 --no-stamp

  # Prefix normalized to lowercase ASCII words
  namo -p "Debug Output / API"

  # Prefix passed through as given
  namo --raw-prefix Build_42`,
		Version:            Version,
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
		RunE: run(func(cmd *cobra.Command, _ []string) error {
			return runRoot(cmd, flags)
		}),
	}

	f := root.Flags()
	f.StringVarP(&flags.prefix, "prefix", "p", "", "prefix, normalized to lowercase ASCII words")
	f.StringVar(&flags.rawPrefix, "raw-prefix", "", "prefix passed through as given (trusted input)")
	f.IntVarP(&flags.count, "count", "n", 1, "number of names, 1 to 100")
	f.StringVarP(&flags.size, "size", "s", string(namegen.SizeStandard), "slug size: short|standard|long")
	f.StringVar(&flags.stamp, "stamp", namegen.DefaultStampLayout, "strftime-style stamp layout")
	f.BoolVar(&flags.shortStamp, "short-stamp", false, "use an HHMM stamp (--stamp %H%M)")
	f.BoolVar(&flags.noStamp, "no-stamp", false, "omit the stamp and keep any prefix")
	root.MarkFlagsMutuallyExclusive("stamp", "short-stamp", "no-stamp")
	root.MarkFlagsMutuallyExclusive("prefix", "raw-prefix")
	root.InitDefaultHelpFlag()
	root.InitDefaultVersionFlag()

	root.AddCommand(newPatternCmd(), newStampsCmd(), newExitCodesCmd())

	return root
}

func runRoot(cmd *cobra.Command, flags *rootFlags) error {
	if err := namegen.ValidateCount(flags.count); err != nil {
		return err
	}

	prefix := flags.rawPrefix
	if cmd.Flags().Changed("prefix") {
		var err error
		prefix, err = namegen.NormalizePrefix(flags.prefix)
		if err != nil {
			return err
		}
	}

	layout := flags.stamp
	switch {
	case flags.noStamp:
		layout = ""
	case flags.shortStamp:
		layout = "%H%M"
	}

	size, err := namegen.ParseSize(flags.size)
	if err != nil {
		return err
	}

	names, err := namegen.Generate(namegen.Options{
		Count:  flags.count,
		Prefix: prefix,
		Size:   size,
		Stamp:  layout,
	})
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	for _, name := range names {
		if _, err := fmt.Fprintln(out, name); err != nil {
			return err
		}
	}
	return nil
}
