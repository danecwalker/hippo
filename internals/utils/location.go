package utils

import "fmt"

type Location struct {
	Line     int
	Column   int
	Filename string
}

func NewLocation(line, column int, filename string) Location {
	return Location{line, column, filename}
}

func (l Location) String() string {
	return fmt.Sprintf("%s:%d:%d", l.Filename, l.Line, l.Column)
}
