package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/200sc/foptgen"
)

// foptgen generates functional options to populate a struct

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

const version = "0.0.2"

var help = flag.Bool("help", false, "display this help text")
var typ = flag.String("struct", "", "struct type to create functional options for (required)")
var output = flag.String("o", "opts_gen.go", "file name to write to")
var genPackage = flag.String("package", "", "go package the written file should belong to. Defaults to what is found in the target directory")
var directory = flag.String("dir", ".", "target directory to search for struct definition")
var fiximports = flag.Bool("fix-imports", true, "run 'goimports' after generation")
var outputToDirectory = flag.Bool("output-to-dir", true, "prepend output with --dir")
var printVersion = flag.Bool("version", false, "print version and exit")
var overwrite = flag.Bool("overwrite", true, "overwrite existing same-named output file")
var optionTypeName = flag.String("option-type-name", "Option", "the functional option type name")

func run() error {
	flag.Parse()
	if *help {
		fmt.Println("foptgen generates functional option declarations to populate a struct")
		fmt.Println()
		fmt.Println("ex: foptgen --struct=Constructor --dir=./component/target")
		fmt.Println()
		flag.PrintDefaults()
		return nil
	}
	if *printVersion {
		fmt.Println("foptgen", version)
		return nil
	}
	if *typ == "" {
		flag.PrintDefaults()
		return fmt.Errorf("required flag %q was not provided", "struct")
	}
	if *outputToDirectory {
		*output = filepath.Join(*directory, *output)
	}

	if !(*overwrite) {
		if err := foptgen.CheckIfOutputExists(*output); err != nil {
			return err
		}
	}

	structDef, fset, pkg, err := foptgen.FindStructInDirectory(*directory, *typ)
	if err != nil {
		return err
	}
	if *genPackage == "" {
		*genPackage = pkg
	}

	tplInput := foptgen.NewTemplateInput(structDef, fset, *typ, *genPackage, *optionTypeName)

	outFile, err := os.Create(*output)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer outFile.Close()

	if err := foptgen.WriteTemplate(outFile, tplInput); err != nil {
		return err
	}

	if *fiximports {
		if err := exec.Command("goimports", "-w", *output).Run(); err != nil {
			return fmt.Errorf("failed to run goimports: %w", err)
		}
	}

	return nil
}
