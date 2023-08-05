package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/danecwalker/hippo/internals/compiler"
)

func main() {
	progam_name := filepath.Base(os.Args[0])
	_ = progam_name
	if len(os.Args[1:]) < 1 {
		fmt.Fprintln(os.Stderr, "No input file")
		os.Exit(1)
	}

	file_name := os.Args[1]

	compiler.Compile(file_name)
}
