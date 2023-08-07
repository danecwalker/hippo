package syntax

import "bytes"

type FuncType struct {
	Params  []*NField
	Results []*Identifier
}

func NewFuncType(params []*NField, results []*Identifier) *FuncType {
	return &FuncType{
		Params:  params,
		Results: results,
	}
}

func (ft *FuncType) expressionNode() {}
func (ft *FuncType) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString("FuncType:\n")
	addIndent(w, indent+1)
	w.WriteString("Params:\n")
	for _, param := range ft.Params {
		param.PrettyPrint(w, indent+2)
	}

	addIndent(w, indent+1)
	w.WriteString("Results:\n")
	for _, result := range ft.Results {
		result.PrettyPrint(w, indent+2)
	}
}
