// Package cli wires the cobra commands for the namo binary.
package cli

import (
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
	version    bool
}

// Version is set at build time via -ldflags.
var Version = "n/a"

// The command groups that root help lists: the commands namo runs, then the
// help topics.
const (
	commandGroupID   = "commands"
	helpTopicGroupID = "topics"
)

// Execute parses the command line, runs the requested command, and returns
// the process exit status: 0 on success, 2 on a usage error the command
// library rejects before the invocation is accepted, and 1 on every other
// failure. A failure prints one line, "namo: <message>", to stderr.
func Execute() int {
	return execute(newRootCmd())
}

// execute runs root with the arguments and streams set on it and returns the
// exit status. The root's persistent pre-run hook, or the help function for
// --help, records that the invocation was accepted; an error before that is
// a usage error and an error after it, from any runner including the
// library's own, is a command failure. A line break or carriage return
// inside a message is shown as \n or \r so the diagnostic stays one line.
func execute(root *cobra.Command) int {
	accepted := false
	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// The root has subcommands, so it leaves Args unset and Cobra
		// accepts any root operand that follows '--'; reject it here.
		if !cmd.HasParent() {
			if err := cobra.NoArgs(cmd, args); err != nil {
				return err
			}
		}
		// Cobra checks required flags and flag groups after this hook, so
		// check them here to keep those failures usage errors.
		if err := cmd.ValidateRequiredFlags(); err != nil {
			return err
		}
		if err := cmd.ValidateFlagGroups(); err != nil {
			return err
		}
		accepted = true
		return nil
	}

	// Cobra prints help for --help before it checks operands and flag
	// groups, and when the help cannot be written it prints the bare error
	// and exits 0. Check them first, then write the help so a failure
	// becomes the one failure line below. The help command and the help
	// topics print help without --help and return their own write errors.
	var helpErr error
	printHelp := root.HelpFunc()
	root.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if help, _ := cmd.Flags().GetBool("help"); !help {
			printHelp(cmd, args)
			return
		}
		if helpErr = checkHelpInvocation(cmd); helpErr != nil {
			return
		}
		accepted = true
		helpErr = writeHelp(cmd, printHelp)
	})

	err := root.Execute()
	if err == nil {
		err = helpErr
	}
	if err == nil {
		return 0
	}
	message := strings.NewReplacer("\r", `\r`, "\n", `\n`).Replace(err.Error())
	_, _ = fmt.Fprintf(root.ErrOrStderr(), "namo: %s\n", message)
	if accepted {
		return 1
	}
	return 2
}

// checkHelpInvocation returns the usage error cmd reports without --help for
// the operands and flag groups given with it. The root accepts no operands,
// though it leaves Args unset.
func checkHelpInvocation(cmd *cobra.Command) error {
	if operands := cmd.Flags().Args(); len(operands) > 0 {
		var err error
		if cmd.HasParent() {
			err = cmd.ValidateArgs(operands)
		} else {
			err = cobra.NoArgs(cmd, operands)
		}
		if err != nil {
			return err
		}
	}
	return cmd.ValidateFlagGroups()
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

namo has no configuration file, writes no files, and reads nothing from
stdin; every option is a flag. The stamp follows the local time zone,
read from /etc/localtime when TZ is unset, or from the zone file TZ
names: an absolute path, or a name such as Asia/Tokyo looked up under
/usr/share/zoneinfo, /usr/share/lib/zoneinfo, /usr/lib/locale/TZ,
/etc/zoneinfo, and, when GOROOT is set, $GOROOT/lib/time/zoneinfo.zip.
An empty, UTC, or unknown TZ, or a missing /etc/localtime, gives UTC.
--no-stamp reads no zone data.

Names go to stdout, one per line, and that plain text is the only output
format; help and version output go to stdout too. namo never prompts,
works without a terminal, runs no external program, and never accesses
the network. Without --raw-prefix or a custom --stamp, every name
consists only of lowercase ASCII letters, digits, and hyphens, so it is
safe inside "$(namo ...)" and in file paths. --raw-prefix and --stamp
are trusted input: line breaks, control bytes, and path separators in
them pass through unchanged.

` + failureContract,
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
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runRoot(cmd, flags)
		},
	}

	f := root.Flags()
	f.StringVarP(&flags.prefix, "prefix", "p", "", "prefix, normalized to lowercase ASCII words")
	f.StringVar(&flags.rawPrefix, "raw-prefix", "", "prefix passed through as given (trusted input)")
	f.IntVarP(&flags.count, "count", "n", 1, "number of names, 1 to 100")
	f.StringVarP(&flags.size, "size", "s", string(namegen.SizeStandard), "slug size: short|standard|long")
	f.StringVar(&flags.stamp, "stamp", namegen.DefaultStampLayout, "strftime-style stamp layout")
	f.BoolVar(&flags.shortStamp, "short-stamp", false, "use an HHMM stamp (--stamp %H%M)")
	f.BoolVar(&flags.noStamp, "no-stamp", false, "omit the stamp and keep any prefix")
	f.BoolVarP(&flags.version, "version", "v", false, "print the version")
	// Shell completion offers nothing after --help or --version only when
	// the flag carries this annotation, which Cobra sets on its own flags.
	_ = f.SetAnnotation("version", cobra.FlagSetByCobraAnnotation, []string{"true"})
	root.MarkFlagsMutuallyExclusive("stamp", "short-stamp", "no-stamp")
	root.MarkFlagsMutuallyExclusive("prefix", "raw-prefix")
	root.InitDefaultHelpFlag()

	root.AddGroup(
		&cobra.Group{ID: commandGroupID, Title: "Commands:"},
		&cobra.Group{ID: helpTopicGroupID, Title: "Help topics:"},
	)
	// No namo command takes a file operand, so shell completion offers no
	// file names unless a command offers its own completions.
	root.CompletionOptions.SetDefaultShellCompDirective(cobra.ShellCompDirectiveNoFileComp)
	help := newHelpCmd()
	root.SetHelpCommand(help)
	root.AddCommand(help, newCompletionCmd(), newPatternCmd(), newStampsCmd(), newExitCodesCmd())

	return root
}

func runRoot(cmd *cobra.Command, flags *rootFlags) error {
	if flags.version {
		_, err := fmt.Fprintf(cmd.OutOrStdout(), "namo version %s\n", Version)
		return err
	}

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
