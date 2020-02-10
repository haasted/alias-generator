package cmd

import (
	"fmt"
	"go/format"
	"path/filepath"
	"strings"
)

func writeAliasFile(physicaldir string, typemap map[string]packageDeclarations) {
	// Clear out package paths that do not have any content to alias.
	for k, v := range typemap {
		if len(v.types) == 0 && len(v.functions) == 0 && len(v.consts) == 0 && len(v.variables) == 0 {
			delete(typemap, k)
		}
	}

	// TODO Get rid of root package

	sb := new(strings.Builder)

	fmt.Fprintf(sb, "package %v\n\n", filepath.Base(physicaldir))

	fmt.Fprintln(sb, "import(")
	for k, _ := range typemap {
		fmt.Fprintf(sb, "\"%v\"\n", k)
	}
	fmt.Fprintln(sb, ")")

	fmt.Fprintln(sb, "const(")
	for pkg, types := range typemap {
		packageAlias := filepath.Base(pkg)
		for _, c := range types.consts {
			fmt.Fprintf(sb, "%v = %v.%v\n", c, packageAlias, c)
		}
	}

	fmt.Fprintln(sb, ")")

	fmt.Fprintln(sb, "var (")
	fmt.Fprintln(sb, "// functions aliases")
	for pkg, types := range typemap {
		packageAlias := filepath.Base(pkg)
		for _, c := range types.functions {
			fmt.Fprintf(sb, "%v = %v.%v\n", c, packageAlias, c)
		}
	}

	fmt.Fprintln(sb, "\n\n// variable aliases")
	for pkg, types := range typemap {
		packageAlias := filepath.Base(pkg)
		for _, c := range types.variables {
			fmt.Fprintf(sb, "%v = %v.%v\n", c, packageAlias, c)
		}
	}

	fmt.Fprintln(sb, ")") // End var

	fmt.Fprintln(sb, "type (")
	for pkg, types := range typemap {
		packageAlias := filepath.Base(pkg)
		for _, c := range types.types {
			fmt.Fprintf(sb, "%v = %v.%v\n", c, packageAlias, c)
		}
	}
	fmt.Fprintln(sb, ")") // End type

	fmtsrc, err := format.Source([]byte(sb.String()))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(fmtsrc))
}
