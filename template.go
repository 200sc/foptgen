package foptsgen

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"html/template"
	"io"
)

//go:embed templates/opts_gen.go.tpl
var TemplateFile string

type TemplateInput struct {
	PackageName string
	Options     []TemplateOption
	StructName  string
}

type TemplateOption struct {
	// struct fields can have multiple names but we're ignoring that for now
	FieldName string
	FieldType string
}

func WriteTemplate(w io.Writer, tplInput *TemplateInput) error {
	tpl, err := template.New("foptgen").Parse(TemplateFile)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	err = tpl.Execute(w, &tplInput)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

func NewTemplateInput(structDef *ast.StructType, fset *token.FileSet, structName, pkg string) *TemplateInput {
	tplInput := &TemplateInput{
		PackageName: pkg,
		StructName:  structName,
	}
	for _, fd := range structDef.Fields.List {
		buff := bytes.NewBuffer([]byte{})
		// Assumption: these are valid ast nodes, so this should not error
		_ = printer.Fprint(buff, fset, fd.Type)
		tplInput.Options = append(tplInput.Options, TemplateOption{
			FieldName: fd.Names[0].Name,
			FieldType: buff.String(),
		})
	}
	return tplInput
}
