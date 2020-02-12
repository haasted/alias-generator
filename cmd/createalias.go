package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"
)

var CreateAlias = &cobra.Command{
	Use:  "create",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Verify input
		for _, dir := range args {
			_, err := ioutil.ReadDir(dir)
			if err != nil {
				return err
			}
		}

		return create(args)
	},
}

func create(dirs []string) error {
	for _, dir := range dirs {
		fullPackage, err := determineFullPackage(dir)
		if err != nil {
			return err
		}

		typesMap := make(map[string]packageDeclarations)
		scanSubdirectories(dir, fullPackage, typesMap)

		// Do not alias the root package
		delete(typesMap, fullPackage)

		writeAliasFile(dir, typesMap)
	}

	return nil
}

// Traverse up until at go.mod file is found. Return the absolute path of the package in the directory provided.
func determineFullPackage(dir string) (string, error) {
	fullPath, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	subPath := make([]string, 0)

	// TODO Wonder what this looks like on Windows OS
	for fullPath != "/" {
		subPath = append([]string{filepath.Base(fullPath)}, subPath...)
		fullPath = filepath.Dir(fullPath)

		modulename, err := parseGoMod(fullPath)
		if err == nil {
			fullPackagePath := fmt.Sprintf("%v/%v", modulename, strings.Join(subPath, "/"))
			return fullPackagePath, nil
		}
	}

	return "", errors.New(fmt.Sprintf("unable to locate go.mod from path %v", dir))
}

// Parse go.mod-file and return root module path.
func parseGoMod(rootProjectPath string) (string, error) {
	bz, err := ioutil.ReadFile(path.Join(rootProjectPath, "go.mod"))
	if err != nil {
		return "", err
	}

	mf, err := modfile.ParseLax("go.mod", bz, nil)
	if err != nil {
		return "", err
	}

	return mf.Module.Mod.Path, nil
}

func filterGoFiles(fi os.FileInfo) bool {
	if strings.HasSuffix(fi.Name(), "_test.go") {
		return false
	}

	if !strings.HasSuffix(fi.Name(), ".go") {
		return false
	}

	if fi.Name() == "alias.go" {
		return false
	}

	return true
}

func startsWithUppercase(s string) bool {
	// TODO There MUST be a more elegant way to check if a string starts with a capital letter
	for _, r := range s {
		return unicode.IsUpper(r)
	}

	return false
}
