package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/danecwalker/hippo/internal/parse"
)

func main() {
	progam_name := filepath.Base(os.Args[0])
	_ = progam_name
	if len(os.Args[1:]) < 1 {
		fmt.Fprintln(os.Stderr, "No input file")
		os.Exit(1)
	}

	file_name := os.Args[1]

	fmt.Println("Parsing file: ", file_name)

	par := parse.NewParser(file_name)

	par.ParseProgram()
}
