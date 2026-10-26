package cli

import "github.com/spf13/cobra"

func newPatternCmd() *cobra.Command {
	return newHelpTopic("pattern", "The [prefix-][stamp-]slug name pattern and slug sizes",
		`Every name namo prints has the form [prefix-][stamp-]slug:

  debug-output-260711154501-star-studded-booze-cruise
  |            |            |
  prefix       stamp        slug

prefix  An optional label you supply, so names group by purpose. A
        normalized prefix is lowercase ASCII letters and digits joined
        by single hyphens; a raw prefix is passed through as given.
stamp   The local time at which the batch was generated, yymmddhhmmss
        by default, so names sort by creation time. Every name in a
        batch shares one stamp. 'namo help stamps' describes layouts.
slug    Random words from hotdiva2000, lowercase ASCII letters and
        digits joined by single hyphens, so a name is easy to recognize
        in a directory listing or a log. No slug repeats within a batch;
        separate runs may repeat one.

Slug sizes:
  short     a modifier and a noun, e.g. undercover-bath-salt
  standard  short, with a word sometimes added at either end, e.g.
            semi-loyal-nail-technician or not-a-venetian-mime-type
  long      a word always added at each end, e.g.
            his-majestys-for-pay-game-dev-part-2

'namo --help' lists the flags that set each part and states when a name
is safe in command substitution and file paths.`)
}
