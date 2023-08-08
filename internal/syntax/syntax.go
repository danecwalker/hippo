package syntax

import "bytes"

type Node interface {
	PrettyPrint(w *bytes.Buffer, indent int)
	Pos() *Position
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}
