package logger

import (
	"fmt"
	"os"

	"github.com/danecwalker/hippo/internal/tok"
)

type LogCode int

const (
	_        LogCode = iota
	LogError         // Error
	LogFatal         // Fatal error
	LogWarn          // Warning

	info_logs

	LogInfo // Information
)

func (e LogCode) String() string {
	switch e {
	case LogError:
		return "\033[31m[ERROR]\033[0m"
	case LogFatal:
		return "\033[31m[ERROR]\033[0m"
	case LogWarn:
		return "\033[33m[WARN]\033[0m"
	case LogInfo:
		return "\033[34m[INFO]\033[0m"
	default:
		panic("invalid issue code")
	}
}

type Log struct {
	Code LogCode
	Msg  string
	Pos  *tok.Position
}

func (e Log) String() string {
	err := ""
	if e.Pos != nil {
		err += fmt.Sprintf("%s:%d:%d\n\t%s %s\n", e.Pos.Filename, e.Pos.Line, e.Pos.Column, e.Code, e.Msg)
	} else {
		err += fmt.Sprintf("%s %s\n", e.Code, e.Msg)
	}

	return err
}

type LogHandler struct {
	logs []Log
}

func NewLogHandler() *LogHandler {
	return &LogHandler{
		logs: make([]Log, 0),
	}
}

func (e *LogHandler) Log(code LogCode, format string, args ...any) {
	e.logs = append(e.logs, Log{
		Code: code,
		Msg:  fmt.Sprintf(format, args...),
		Pos:  nil,
	})

	if code == LogFatal {
		e.Fatal()
	}
}

func (e *LogHandler) InfoLog(format string, args ...any) {
	log := Log{
		Code: LogInfo,
		Msg:  fmt.Sprintf(format, args...),
		Pos:  nil,
	}

	fmt.Fprintf(os.Stdout, "%s", log)
}

func (e *LogHandler) LogAt(code LogCode, pos *tok.Position, format string, args ...any) {
	e.logs = append(e.logs, Log{
		Code: code,
		Msg:  fmt.Sprintf(format, args...),
		Pos:  pos,
	})

	if code == LogFatal {
		e.Fatal()
	}
}

func (e *LogHandler) MaybeFatal() {
	if len(e.logs) > 0 {
		e.Fatal()
	}
}

func (e *LogHandler) Fatal() {
	for _, log := range e.logs {
		if log.Code < info_logs {
			fmt.Fprintf(os.Stderr, "%s", log)
		} else {
			fmt.Fprintf(os.Stdout, "%s", log)
		}
	}

	os.Exit(1)
}
