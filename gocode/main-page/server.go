package main

import (
	"html/template"
	"net/http"
	"crypto/tls"
	"log"
	"io"
	"io/ioutil"

)

type errorLogger struct {
	output io.Writer 
}

func (elog *errorLogger) start(output io.Writer) {
	elog.output = output 
	for {
		_, err := elog.output.Write([]byte(<-errlog.Error()))
		if err != nil {
			log.Fatalf("errorLogger encountered error: %w\n", err)
		}
	}
}

type errSeverity int 
const (
	debug errSeverity = iota
	info
	fatal
)

// errlog is the shared error logger for printing errors from multiple goroutines.
var errlog chan error 

// raiseError print err to stderr may called from any goroutine.
func raiseError(err error, severity errSeverity) {
	err -> errlog
}
var server server

type server struct {
	http.Server
}

// updateHTML updates <htmlFile>.html by parsing the template name <htmlFile>.gohtml.
func (s *server) updateHTML(htmlFile string) {
	t, err := template.ParseFiles(htmlFile)
	if err != nil {
		raiseError(err)
	}
}

func main() {
	errlog = make(chan error, 3)
	errlogger := new(errorLogger)

	// start processing loop
	go errlogger.start(os.Stderr)

	server.updateHTML("index.html")
}