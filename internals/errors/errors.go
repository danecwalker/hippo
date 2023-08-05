package errors

import (
	"fmt"
	"os"
)

type ErrorLevel int

const (
	ERROR ErrorLevel = iota
	WARN
)

func (el ErrorLevel) String() string {
	switch el {
	case ERROR:
		// Red terminal color code
		return "\033[31m[ERROR]\033[0m"
	case WARN:
		// Yellow terminal color code
		return "\033[33m[WARN]\033[0m"
	default:
		panic(fmt.Sprintf("Undefined string representation for ErrorLevel %d", el))
	}
}

type Error struct {
	level ErrorLevel
	msg   string
	loc   string
}

func (e Error) String() string {
	if e.loc != "" {
		return fmt.Sprintf("%s\n\t%s\t%s", e.loc, e.level, e.msg)
	}

	return fmt.Sprintf("%s %s", e.level, e.msg)
}

type ErrorHandler struct {
	errors []*Error
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		errors: make([]*Error, 0),
	}
}

type OptFunc func(*Error) *Error

func (eh *ErrorHandler) ErrorWithLoc(level ErrorLevel, fmt_string string, loc string, args ...any) {
	err := &Error{
		level: level,
		msg:   fmt.Sprintf(fmt_string, args...),
		loc:   loc,
	}

	eh.errors = append(eh.errors, err)
}

func (eh *ErrorHandler) Error(level ErrorLevel, fmt_string string, args ...any) {
	err := &Error{
		level: level,
		msg:   fmt.Sprintf(fmt_string, args...),
		loc:   "",
	}

	eh.errors = append(eh.errors, err)
}

func (eh *ErrorHandler) ShouldExit() {
	if len(eh.errors) > 0 {
		for _, e := range eh.errors {
			fmt.Fprintf(os.Stderr, "%s\n", e)
		}
		os.Exit(1)
	}
}
