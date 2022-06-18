package logger

import (
	"fmt"
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

func logf_(level logLevel, msg string, args ...interface{}) {
	log_(fmt.Sprintf(msg, args...), level)
}

func Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		logf_(_INFO, msg, args...)
	} else {
		log_(msg, _INFO)
	}
}
func Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		logf_(_ERROR, msg, args...)
	} else {
		log_(msg, _ERROR)
	}
}
