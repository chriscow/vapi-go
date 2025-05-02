package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "vapi",
		Short: "VAPI CLI utility",
		Long:  "A command line utility for interacting with the VAPI service",
	}

	rootCmd.AddCommand(newCallCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
