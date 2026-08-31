// Package cmd wires the cobra commands for the namo binary.
package cmd

import (
	"fmt"

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

// Execute parses the command line and runs the requested command.
func Execute() error {
	return newRootCmd().Execute()
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
-s/--size short|standard|long sets slug length: short is a modifier and a
noun, standard sometimes adds a word at either end, long always adds one
at each end. --stamp LAYOUT renders a strftime-style layout with %Y %y %m
%d %H %M %S and %%, --short-stamp is --stamp %H%M, and --no-stamp omits
the stamp and keeps any prefix. --prefix and --raw-prefix are mutually
exclusive, as are --stamp, --short-stamp, and --no-stamp.

namo has no configuration file, reads nothing from stdin, and takes every
option as a flag. The stamp follows the local time zone.

Names go to stdout, one per line; help and version output go to stdout
too. Errors go to stderr as 'namo: error: <message>' with no usage text.
namo never prompts, runs no external program, and never accesses the
network. Without --raw-prefix or a custom --stamp, every name consists
only of lowercase ASCII letters, digits, and hyphens, so it is safe inside
"$(namo ...)" and in file paths. --raw-prefix and --stamp are trusted
input: line breaks, control bytes, and path separators in them pass
through unchanged.

Run 'namo docs' for stamp layout verbs, size details, and recipes, and
'namo help exit-codes' for exit status meanings.`,
		Version:       Version,
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runRoot(cmd, flags)
		},
	}

	f := root.Flags()
	f.StringVarP(&flags.prefix, "prefix", "p", "", "ASCII-normalized prefix joined by a hyphen; mutually exclusive with --raw-prefix")
	f.StringVar(&flags.rawPrefix, "raw-prefix", "", "unsafe trusted prefix: non-empty input is preserved before a joining hyphen; empty omits the prefix; mutually exclusive with --prefix")
	f.IntVarP(&flags.count, "count", "n", 1, "number of names to generate (1 <= count <= 100; one timestamp, unique slugs)")
	f.StringVarP(&flags.size, "size", "s", string(namegen.SizeStandard), "slug size: short, standard, or long")
	f.StringVar(&flags.stamp, "stamp", namegen.DefaultStampLayout, "trusted custom strftime layout; preserved literals can make output unsafe (%Y %y %m %d %H %M %S)")
	f.BoolVar(&flags.shortStamp, "short-stamp", false, "use an HHMM timestamp for ephemeral names (--stamp %H%M)")
	f.BoolVar(&flags.noStamp, "no-stamp", false, "omit the timestamp while retaining any prefix")
	root.MarkFlagsMutuallyExclusive("stamp", "short-stamp", "no-stamp")
	root.MarkFlagsMutuallyExclusive("prefix", "raw-prefix")

	root.AddCommand(newDocsCmd(), newExitCodesCmd())

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
