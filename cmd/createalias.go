package cmd

import (
	"errors"
	"fmt"
	"golang.org/x/mod/modfile"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"
)

func startGenerator(dirs []string, update bool) error {
	for _, dir := range dirs {
		fullPackage, err := determineFullPackage(dir)
		if err != nil {
			return err
		}

		typesMap := make(map[string]packageDeclarations)
		scanSubdirectories(dir, fullPackage, typesMap)

		// Do not alias the root package
		delete(typesMap, fullPackage)

		if update {
			aliaspackages, err := parseImports(dir)
			if err != nil {
				return err
			}

			// Create a new typesmap containing only the packages explicitly mentioned in alias.go
			newTypesMap := make(map[string]packageDeclarations)
			for _, pack := range aliaspackages {
				newTypesMap[pack] = typesMap[pack]
			}
			typesMap = newTypesMap
		}

		bz := generateAliasFile(dir, typesMap)
		ioutil.WriteFile(filepath.Join(dir, "alias.go"), bz, os.ModePerm)
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
