package main

import (
	"github.com/haasted/alias-generator/cmd"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "aliasgen",
		Short: "creates or updates alias.go files",
	}
	rootCmd.AddCommand(cmd.CreateCommand(), cmd.UpdateCommand())
	rootCmd.SilenceUsage = true

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}
