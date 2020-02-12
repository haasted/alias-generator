package main

import "github.com/spf13/cobra"
import "github.com/haasted/alias-generator/cmd"

func main() {
	rootCmd := &cobra.Command{
		Use:   "aliasgen",
		Short: "creates or updates alias.go files",
	}
	rootCmd.AddCommand(cmd.CreateAlias)
	rootCmd.Execute()
}
