package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	TESTDIR               = "../testdata"
	TESTDATAPACKAGE       = "github.com/haasted/alias-generator/testdata"
	TESTDATAPACKAGEDB     = "github.com/haasted/alias-generator/testdata/db"
	TESTDATAPACKAGEIMPL   = "github.com/haasted/alias-generator/testdata/impl"
	TESTDATAPACKAGESECRET = "github.com/haasted/alias-generator/testdata/secret"
)

func TestDetermineFullPackage(t *testing.T) {
	fullPackage, err := determineFullPackage(TESTDIR)
	require.NoError(t, err)
	require.Equal(t, TESTDATAPACKAGE, fullPackage)
}

func TestTypeScanner(t *testing.T) {
	typesMap := make(map[string]packageDeclarations)
	scanSubdirectories(TESTDIR, TESTDATAPACKAGE, typesMap)

	require.Len(t, typesMap, 4)

	{
		decls := typesMap[TESTDATAPACKAGEIMPL]
		require.Contains(t, decls.variables, "ThisIsAVar")
		require.Contains(t, decls.variables, "AnotherVar")
		require.NotContains(t, decls.variables, "MoreVars")
		require.NotContains(t, decls.variables, "andAnother")
	}

	{
		decls := typesMap[TESTDATAPACKAGEIMPL]
		allIdentifiers := append(decls.consts, decls.functions...)
		allIdentifiers = append(allIdentifiers, decls.variables...)
		allIdentifiers = append(allIdentifiers, decls.types...)

		for _, id := range allIdentifiers {
			require.NotEqual(t, "_", id)
		}
	}

	{
		// The entire /db package is noaliased in various ways. Ensure nothing gets picked up.
		decls := typesMap[TESTDATAPACKAGEDB]

		require.Empty(t, decls.types)
		require.Empty(t, decls.variables)
		require.Empty(t, decls.functions)
		require.Empty(t, decls.consts)
	}

	{
		// The entire /secret package is noaliased by annotating the package declaration.
		decls := typesMap[TESTDATAPACKAGESECRET]

		require.Empty(t, decls.types)
		require.Empty(t, decls.variables)
		require.Empty(t, decls.functions)
		require.Empty(t, decls.consts)
	}
}
