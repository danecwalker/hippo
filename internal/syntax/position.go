package syntax

import "fmt"

type Position struct {
	Offset   int
	Line     int
	Column   int
	Filename string
}

func NewPosition(offset int, line int, column int, filename string) *Position {
	return &Position{
		Offset:   offset,
		Line:     line,
		Column:   column,
		Filename: filename,
	}
}

func (p *Position) String() string {
	return fmt.Sprintf("%s:%d:%d", p.Filename, p.Line, p.Column)
}
