package main

import (
	"fmt"
	"io"
	"log"
)

type errSeverity int

const (
	NONE              = -1
	DEBUG errSeverity = iota
	ERROR
	FATAL
	INFO
)

// taggedErr implements the error interface, and holds severity info for the logger as well as the error.
type taggedErr struct {
	err      error
	severity errSeverity
}

func tagErr(err error, severity errSeverity) taggedErr {
	return taggedErr{err, severity}
}

// Implements error interface.
func (te *taggedErr) Error() string {
	return te.err.Error()
}

// errlog is the shared error logger for printing errors from multiple goroutines.
var errlog chan taggedErr

// throw print err to stderr may called from any goroutine.
// If severity is nil, ERROR is assumed.
func throw(err error, severity ...errSeverity) {
	var tag taggedErr
	if severity != nil {
		tag = tagErr(err, severity[0])
	}
	tag = tagErr(err, ERROR)
	errlog <- tag
}

type errorLogger struct {
	out io.Writer
}

func (elog *errorLogger) run(out io.Writer) {
	elog.out = out
	for {
		_, err := elog.out.Write([]byte(elog.Error()))
		if err != nil {
			log.Fatalf(fmt.Errorf("errorLogger encountered error: %w\n", err).Error())
		}
	}
}

func (elog *errorLogger) Error() string {
	e, ok := <-errlog
	if ok {
		return e.Error()
	}
	return "" // TODO: check how std lib handles this
}
