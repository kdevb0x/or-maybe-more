package app

import (
	"io"
	"log"
)

type logger struct {
	out *log.Logger
}

func newLogger(output io.Writer) *logger {
	return &logger{log.New(output, "oom-server: ", log.LstdFlags|log.Llongfile)}
}

func (l *logger) Println(s string) {
	return l.out.Println(s)
}

func (l *logger) Printf(format string, v ...interface{}) {
	return l.out.Printf(format, v)
}
