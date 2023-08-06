package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/danecwalker/hippo/internal/flow"
	"github.com/danecwalker/hippo/internal/par"
)

func main() {
	progam_name := filepath.Base(os.Args[0])
	_ = progam_name
	if len(os.Args[1:]) < 1 {
		fmt.Fprintln(os.Stderr, "No input file")
		os.Exit(1)
	}

	file_name := os.Args[1]

	p := par.NewParser(file_name)
	f := p.Parse()

	b, _ := json.MarshalIndent(f, "", " ")
	fmt.Println(string(b))

	cfg := flow.NewCFG()
	cfg.GenCFG(f)
}
