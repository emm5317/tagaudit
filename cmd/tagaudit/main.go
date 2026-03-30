package main

import (
	"fmt"
	"os"

	"github.com/emm5317/tagaudit/cmd/tagaudit/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
