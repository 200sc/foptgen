package foptgen

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// FindStructInDirectory reads a directory, without recursing to subdirectories. It searches all
// Go files within the directory for a struct named 'structName'. It will return that struct's definition,
// the token file set used to parse that file, and the name of the package that struct was found in, or
// an error.
func FindStructInDirectory(directory, structName string) (*ast.StructType, *token.FileSet, string, error) {
	dirEntries, err := os.ReadDir(directory)
	if err != nil {
		return nil, nil, "", fmt.Errorf("error reading directory: %w", err)
	}
	for _, ent := range dirEntries {
		if ent.IsDir() {
			continue
		}
		if !strings.HasSuffix(ent.Name(), ".go") {
			continue
		}
		structDef, fset, pkg, err := FindStructInFile(filepath.Join(directory, ent.Name()), structName)
		if err != nil && !errors.Is(err, ErrStructNotFoundInFile) {
			return nil, nil, "", fmt.Errorf("error reading %v: %w", ent.Name(), err)
		}
		if structDef != nil {
			return structDef, fset, pkg, nil
		}
	}
	return nil, nil, "", fmt.Errorf("struct %q not found", structName)
}

// FindStructInFile searches a fully qualified file path for a struct named 'structName'.
// If found, it will return that struct definition, the token file set used to parse the file,
// and the package name of the file. Otherwise, if not found, it will return ErrStructNotFoundInFile.
func FindStructInFile(file, structName string) (structDef *ast.StructType, fset *token.FileSet, pkg string, err error) {
	fset = token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to parse go file: %w", err)
	}

	var lastIdent *ast.Ident
	pkg = f.Name.Name

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			lastIdent = x
		case *ast.StructType:
			if lastIdent.Name == structName {
				structDef = x
				return false
			}
		}
		return true
	})
	if structDef == nil {
		err = ErrStructNotFoundInFile
	}

	return
}

var ErrStructNotFoundInFile = fmt.Errorf("struct not found in file")
