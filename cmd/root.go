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
		Long: "namo composes [prefix-]stamp-slug names: a sortable timestamp plus a\n" +
			"memorable random slug.\n\n" +
			"  namo -p debug-output  ->  debug-output-260711154501-star-studded-booze-cruise",
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
	f.StringVar(&flags.rawPrefix, "raw-prefix", "", "unsafe trusted prefix preserved before a joining hyphen; mutually exclusive with --prefix")
	f.IntVarP(&flags.count, "count", "n", 1, "number of names to generate (one timestamp, unique slugs)")
	f.StringVarP(&flags.size, "size", "s", string(namegen.SizeStandard), "slug size: short, standard, or long")
	f.StringVar(&flags.stamp, "stamp", namegen.DefaultStampLayout, "strftime-style timestamp layout (%Y %y %m %d %H %M %S)")
	f.BoolVar(&flags.shortStamp, "short-stamp", false, "use an HHMM timestamp for ephemeral names (--stamp %H%M)")
	f.BoolVar(&flags.noStamp, "no-stamp", false, "omit the timestamp")
	root.MarkFlagsMutuallyExclusive("stamp", "short-stamp", "no-stamp")
	root.MarkFlagsMutuallyExclusive("prefix", "raw-prefix")

	root.AddCommand(newDocsCmd())

	return root
}

func runRoot(cmd *cobra.Command, flags *rootFlags) error {
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
