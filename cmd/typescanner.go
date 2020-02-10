package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path"
)

func scanSubdirectories(physDir, absolutePackage string, types map[string]packageDeclarations) error {
	files, err := ioutil.ReadDir(physDir)
	if err != nil {
		return err
	}

	declarations := packageDeclarations{}
	expr, err := parser.ParseDir(token.NewFileSet(), physDir, filterGoFiles, parser.ParseComments)
	for _, e := range expr {
		for _, f := range e.Files {
			for _, d := range f.Decls {
				switch decl := d.(type) {
				case *ast.GenDecl:
					for _, s := range decl.Specs {
						switch spec := s.(type) {
						case *ast.ImportSpec:

						case *ast.ValueSpec:
							names := make([]string, 0)
							for _, name := range spec.Names {
								if name.Name == "_" {
									continue
								}

								if !startsWithUppercase(name.Name) {
									continue
								}

								names = append(names, name.Name)
							}

							switch decl.Tok {
							case token.CONST:
								declarations.consts = append(declarations.consts, names...)
							case token.VAR:
								declarations.variables = append(declarations.variables, names...)
							default:
								fmt.Println("Unhandled declaration token: ", decl.Tok)
							}

						case *ast.TypeSpec:
							if startsWithUppercase(spec.Name.Name) {
								declarations.types = append(declarations.types, spec.Name.Name)
							}

						default:
							fmt.Printf("Unhandled spec type : %T\n", spec)
						}
					}

				case *ast.FuncDecl:
					if decl.Recv == nil {
						if startsWithUppercase(decl.Name.Name) {
							declarations.functions = append(declarations.functions, decl.Name.Name)
						}
					}

				default:
					fmt.Printf("Unhandled decl type : %T\n", d)
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
