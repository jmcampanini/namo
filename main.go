// Command namo generates memorable, sortable names.
package main

import (
	"os"

	"github.com/jmcampanini/namo/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
