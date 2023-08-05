package compiler

import (
	"fmt"
	"os"

	"github.com/danecwalker/hippo/internals/errors"
	"github.com/danecwalker/hippo/internals/lex"
	"github.com/danecwalker/hippo/internals/parse"
	"github.com/danecwalker/hippo/internals/ssair"
	"github.com/danecwalker/hippo/internals/tree"
)

func Compile(file_name string) {
	mainErrHandle := errors.NewErrorHandler()
	f, err := os.Open(file_name)
	if err != nil {
		mainErrHandle.Error(errors.ERROR, "could not open file %s", file_name)
		mainErrHandle.ShouldExit()
	}

	defer f.Close()

	info, err := os.Stat(file_name)
	if err != nil {
		mainErrHandle.Error(errors.ERROR, "could not open file %s", file_name)
		mainErrHandle.ShouldExit()
	}

	fsize := info.Size()

	buf := make([]byte, fsize)
	_, err = f.Read(buf)
	if err != nil {
		mainErrHandle.Error(errors.ERROR, "could not read file %s", file_name)
		mainErrHandle.ShouldExit()
	}

	l := lex.NewLexer(file_name, buf)
	file := parse.ParseFile(l, mainErrHandle)
	mainErrHandle.ShouldExit()

	// tree.Print(file)
	// fmt.Println()
	// fmt.Println("--------------------------------------------")
	// fmt.Println()
	ur := parse.FindUnresolved(file, mainErrHandle)
	tree.Print(file)
	// fmt.Println()
	// fmt.Println("--------------------------------------------")
	// fmt.Println()
	fmt.Printf("Unresolved items (len = %d)\n", len(ur))
	for _, ident := range ur {
		fmt.Printf("%s @ (%s)\n", ident.Name, ident.NamePos)
	}
	mainErrHandle.ShouldExit()

	parse.Typecheck(file, mainErrHandle)
	mainErrHandle.ShouldExit()

	ssair.GenForm(file, mainErrHandle)

	// fmt.Println(file.Decls[0].(*tree.StoreDecl).Specs[0].(*tree.ValueSpec).Names[0].Obj.Name)
}
