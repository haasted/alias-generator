package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

func scanSubdirectories(physDir, absolutePackage string, types map[string]packageDeclarations) error {
	files, err := ioutil.ReadDir(physDir)
	if err != nil {
		return err
	}

	declarations := packageDeclarations{}
	expr, err := parser.ParseDir(token.NewFileSet(), physDir, filterGoFiles, parser.ParseComments)

stopscan:
	for _, e := range expr {
		for _, f := range e.Files {
			// The entire package can be noaliased by annotating the package declaration
			if markedExempt(f.Doc) {
				declarations = packageDeclarations{}
				break stopscan
			}

			for _, d := range f.Decls {
				switch decl := d.(type) {

				// Capture function declarations
				case *ast.FuncDecl:
					if decl.Recv == nil { // Function, not receiver
						if startsWithUppercase(decl.Name.Name) && !markedExempt(decl.Doc) {
							declarations.functions = append(declarations.functions, decl.Name.Name)
						}
					}

				// Capture var, const and type declarations
				case *ast.GenDecl:
					for _, s := range decl.Specs {
						switch spec := s.(type) {

						case *ast.TypeSpec:
							if startsWithUppercase(spec.Name.Name) && !markedExempt(spec.Comment, spec.Doc, decl.Doc) {
								declarations.types = append(declarations.types, spec.Name.Name)
							}

						case *ast.ValueSpec:
							names := make([]string, 0)
							for _, name := range spec.Names {
								if !startsWithUppercase(name.Name) {
									continue
								}

								if markedExempt(spec.Comment, spec.Doc, decl.Doc) {
									continue
								}

								names = append(names, name.Name)
							}

							switch decl.Tok {
							case token.CONST:
								declarations.consts = append(declarations.consts, names...)
							case token.VAR:
								declarations.variables = append(declarations.variables, names...)
							}
						}
					}

				}
			}
		}
	}

	types[absolutePackage] = declarations

	// Scan sub-directories recursively
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		scanSubdirectories(
			path.Join(physDir, file.Name()),
			fmt.Sprintf("%v/%v", absolutePackage, file.Name()),
			types,
		)
	}

	return nil
}

func parseAliasGenAnnotations(dir string) (aliaspackages []string, err error) {
	aliasfile := filepath.Join(dir, "alias.go")
	file, err := parser.ParseFile(token.NewFileSet(), aliasfile, nil, parser.ParseComments)
	if err != nil {
		return
	}

	if file.Comments == nil {
		return
	}

	for _, c := range file.Comments {
		for _, l := range c.List {
			// Find any comment lines that annotate the packages that should be aliased.
			if strings.Index(l.Text, AliasKeyword) != -1 {
				p := strings.Split(l.Text, AliasKeyword)[1]
				aliaspackages = append(aliaspackages, strings.TrimSpace(p))
			}
		}
	}

	return
}

func markedExempt(cgs ...*ast.CommentGroup) bool {
	for _, cg := range cgs {
		if cg == nil {
			continue
		}

		if strings.Contains(cg.Text(), NoAliasEsc) {
			return true
		}
	}

	return false
}
