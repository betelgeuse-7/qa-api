package logger

import (
	"io"
	"log"
)

type logLevel uint

type logger struct {
	l log.Logger
}

const (
	_INFO logLevel = iota
	_ERROR
)

var l *logger = &logger{l: *log.Default()}

func SetOut(out io.Writer) {
	l.l.SetOutput(out)
}

func log_(msg string, level logLevel) {
	prefix := ""
	switch level {
	case _INFO:
		prefix = "[INFO] "
	case _ERROR:
		prefix = "[ERROR] "
	}
	msg = prefix + " " + msg
	l.l.Println(msg)
}

func Info(msg string)  { log_(msg, _INFO) }
func Error(msg string) { log_(msg, _ERROR) }
