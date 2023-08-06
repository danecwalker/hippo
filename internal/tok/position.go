package tok

import "fmt"

type Position struct {
	Line     int
	Column   int
	Filename string
}

func NewPosition(line int, column int, filename string) Position {
	return Position{
		Line:     line,
		Column:   column,
		Filename: filename,
	}
}

func (p Position) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s:%d:%d\"", p.Filename, p.Line, p.Column)), nil
}
