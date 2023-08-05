package parse

import (
	"fmt"

	"github.com/danecwalker/hippo/internals/symbol"
)

type Precedence int

const (
	_ Precedence = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

func (p Precedence) String() string {
	switch p {
	case LOWEST:
		return "LOWEST"
	case EQUALS:
		return "EQUALS"
	case LESSGREATER:
		return "LESSGREATER"
	case SUM:
		return "SUM"
	case PRODUCT:
		return "PRODUCT"
	case PREFIX:
		return "PREFIX"
	case CALL:
		return "CALL"
	default:
		panic(fmt.Sprintf("Undefined string representation for precedence %d", p))
	}
}

var precedences = map[symbol.SymbolType]Precedence{}
