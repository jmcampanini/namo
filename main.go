// Command namo generates memorable, sortable names.
package main

import (
	"fmt"
	"os"

	"github.com/jmcampanini/namo/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "namo: error:", err)
		os.Exit(1)
	}
}
