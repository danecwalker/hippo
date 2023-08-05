package ssair

import (
	"github.com/danecwalker/hippo/internals/errors"
	"github.com/danecwalker/hippo/internals/tree"
)

type Block struct {
	Id           int
	Parent       *Block
	Children     []*Block
	Instructions []*Instruction
}

type Instruction struct {
	Id   int
	Op   string
	Args []string
}

type Program struct {
	*errors.ErrorHandler
	hasEntry bool
	blocks   []*Block
}

func NewProgram(eh *errors.ErrorHandler) *Program {
	return &Program{ErrorHandler: eh}
}

func (p *Program) Build(ast *tree.File) {
}

func GenForm(ast *tree.File, eh *errors.ErrorHandler) *Program {
	prog := NewProgram(eh)
	prog.Build(ast)
	return prog
}
