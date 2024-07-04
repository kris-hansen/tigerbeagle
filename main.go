package main

import (
	"fmt"
	"os"

	"github.com/kris-hansen/tigerbeagle/internal/app"
	"github.com/kris-hansen/tigerbeagle/internal/cli"
)

func main() {
	tigerBeagle := app.NewTigerBeagle()
	rootCmd := cli.NewRootCommand(tigerBeagle)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
