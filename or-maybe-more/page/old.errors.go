// +build "ignore"

package page

import (
	"fmt"
	"io"
	"log"
)

type ErrSeverity int

const (
	NONE              = -1
	DEBUG ErrSeverity = iota
	ERROR
	FATAL
	INFO
)

// TaggedErr implements the error interface, and holds severity info for the logger as well as the error.
type TaggedErr struct {
	err      error
	severity ErrSeverity
}

func TagErr(err error, severity ErrSeverity) TaggedErr {
	return TaggedErr{err, severity}
}

// Implements error interface.
func (te *TaggedErr) Error() string {
	return te.err.Error()
}

// errlog is the shared error logger for printing errors from multiple goroutines.
var errlog chan TaggedErr

// throw print err to stderr may called from any goroutine.
// If severity is nil, ERROR is assumed.
func throw(err error, severity ...ErrSeverity) {
	var tag TaggedErr
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
