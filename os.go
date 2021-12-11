package foptsgen

import (
	"bufio"
	"fmt"
	"os"
)

// CheckIfOutputExists will prompt to overwrite a file if it already exists. If not confirmed,
// it will error.
func CheckIfOutputExists(filename string) error {
	outFile, err := os.Open(filename)
	if !os.IsNotExist(err) {
		fmt.Printf("Output target %v already exists; overwrite? (y/N)\n", filename)
		scan := bufio.NewScanner(os.Stdin)
		scan.Scan()
		if scan.Text() != "y" {
			return fmt.Errorf("not overwriting existing file")
		}
	}
	outFile.Close()
	return nil
}
