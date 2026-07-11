// Command namo generates memorable, sortable names.
package main

import (
	"fmt"
	"os"

	"github.com/jmcampanini/namo/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "namo: error:", err)
		os.Exit(1)
	}
}
