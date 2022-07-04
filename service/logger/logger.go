package logger

import (
	"fmt"
	"io"
	"log"
)

type logLevel uint

type Logger struct {
	l *log.Logger
}

func NewLogger(l *log.Logger) *Logger {
	return &Logger{l: l}
}

const (
	_INFO logLevel = iota
	_ERROR
)

func (l *Logger) SetOut(out io.Writer) {
	l.l.SetOutput(out)
}

func log_(l *Logger, msg string, level logLevel) {
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

func logf_(l *Logger, level logLevel, msg string, args ...interface{}) {
	log_(l, fmt.Sprintf(msg, args...), level)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		logf_(l, _INFO, msg, args...)
	} else {
		log_(l, msg, _INFO)
	}
}
func (l *Logger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		logf_(l, _ERROR, msg, args...)
	} else {
		log_(l, msg, _ERROR)
	}
}
