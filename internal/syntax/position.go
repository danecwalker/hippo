package syntax

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
