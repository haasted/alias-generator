package cmd

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create [directory...]",
		Short: "Create alias.go files in the specified directories",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Verify input
			for _, dir := range args {
				_, err := ioutil.ReadDir(dir)
				if err != nil {
					return err
				}
			}

			return startGenerator(args)
		},
	}
}

func UpdateCommand() (cmd *cobra.Command) {
	doNotRecurse := false

	cmd = &cobra.Command{
		Use:   "update [directory]",
		Short: "Search for and update existing alias.go files",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := args[0]

			_, err := ioutil.ReadDir(dir)
			if err != nil {
				return err
			}

			aliasFileDirs := make([]string, 0)
			if doNotRecurse {
				aliasFileDirs = []string{dir}
			} else {
				filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
					if filepath.Base(path) == "alias.go" {
						aliasFileDirs = append(aliasFileDirs, filepath.Dir(path))
					}
					return nil
				})
			}

			return startGenerator(aliasFileDirs)
		},
	}

	cmd.Flags().BoolVarP(&doNotRecurse, "no-recurse", "R", false, "Do not traverse sub-directories.")

	return
}
