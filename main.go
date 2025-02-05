package main

import (
	"os"
	"github.com/eriklundjensen/thdctrl/cmd/thdctrl"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "thdctrl",
		Short: "Talos Hetzner Dedicate Servers CLI",
	}

	for _, cmd := range thdctrl.Commands {
			rootCmd.AddCommand(cmd)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
