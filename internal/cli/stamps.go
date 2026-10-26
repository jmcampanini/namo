package cli

import "github.com/spf13/cobra"

func newStampsCmd() *cobra.Command {
	return newHelpTopic("stamps", "Stamp layout verbs and how the stamp is rendered",
		`The stamp is the local time at which namo generated the batch, rendered
through a strftime-style layout. The default layout, %y%m%d%H%M%S, gives
a twelve-digit stamp such as 260711154501 that sorts by creation time.
The short layout, %H%M, gives 1545 for names that only need to be told
apart within a day, and a name can omit the stamp altogether. The flags
that choose a layout are listed in 'namo --help'.

Layout verbs:
  %Y  four-digit year   2026        %H  hour (00-23)         15
  %y  two-digit year    26          %M  minute (00-59)       45
  %m  month (01-12)     07          %S  second (00-59)       01
  %d  day (01-31)       11          %%  literal percent sign

Any other verb is an error, as is a layout that ends with a bare %.
Every other character passes through unchanged, so %Y-%m-%d renders as
2026-07-11. A custom layout is trusted input: line breaks, control
bytes, and path separators in it reach the output as given.`)
}
