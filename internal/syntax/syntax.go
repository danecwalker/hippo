package syntax

import "bytes"

type Node interface {
	PrettyPrint(w *bytes.Buffer, indent int)
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}
