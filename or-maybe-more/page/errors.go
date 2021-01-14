package clipquip

import (
	"errors"
	"fmt"
)

type LogLevel int

const (
	LogLevelFATAL LogLevel = -2
	LogLevelERROR          = -1
	LogLevelINFO           // zero value
	LogLevelDEBUG = 1

	// LogLevels from github.com/jackc/pgx
	//
	// LogLevelTrace = 6
	// LogLevelDebug = 5
	// LogLevelInfo  = 4
	// LogLevelWarn  = 3
	// LogLevelError = 2
	// LogLevelNone  = 1

)

// TaggedErr implements the error interface, and holds level info for the logger as well as the error.
type TaggedErr struct {
	err   error
	level LogLevel
}

func TagErr(err error, level LogLevel) TaggedErr {
	return TaggedErr{err, level}
}

// Implements error interface.
func (te *TaggedErr) Error() string {
	return te.err.Error()
}

// Log implements github.com/jackc/pgx Logger interface.
func (el *ErrLog) Log(level LogLevel, msg string, data map[string]interface{}) {
	if data == nil {
		*el <- TagErr(errors.New(msg), level)

	}
	var tag TaggedErr

	// try parsing placeholders
	if errmsg := fmt.Errorf(msg, data); errmsg.Error() != "" {
		tag = TagErr(errmsg, level)
	}
	*el <- tag

}

// ErrLog is the shared error logger for printing errors from multiple goroutines.
type ErrLog chan TaggedErr

// errlog is the default ErrLog instance.
var errlog ErrLog

// throw print err to stderr may called from any goroutine.
// If level is nil, LogERROR is assumed.
func throw(err error, level LogLevel) {

	// [kdev] why all this again?
	/*
		// a rolling counter for error logging
		var counter uint64

		// counterEpoch is incremented each time counter rolls over.
		// lets hope we never need a "counterEpochEpoch" - [kdev]
		var counterEpoch uint64

		// counter closure
		go func() {
			// if counter is about to overflow
			if counter&atomic.AddUint64(&counter, ^uint64(0)) == 0 {

				// atomically increment the counterEpoch
				counterEpoch = atomic.AddUint64(&counterEpoch, uint64(1))
			}
			counter = atomic.AddUint64(&counter, uint64(1))
		}()
	*/

	var tag TaggedErr
	tag = TagErr(err, level)
	errlog <- tag
}
