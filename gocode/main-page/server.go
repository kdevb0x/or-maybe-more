package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// errlog is the shared error logger for printing errors from multiple goroutines.
var errlog chan error

// raiseError print err to stderr may called from any goroutine.
func raiseError(err error, severity errSeverity) {
	errlog <- err
}

type errorLogger struct {
	out io.Writer
}

func (elog *errorLogger) start(out io.Writer) {
	elog.out = out
	for {
		_, err := elog.out.Write([]byte(<-errlog))
		if err != nil {
			log.Fatalf("errorLogger encountered error: %w\n", err)
		}
	}
}

type errSeverity int

const (
	DEBUG errSeverity = iota
	INFO
	ERROR
	FATAL
)

var Server = new(server)

type server struct {
	http.Server
}

// updateHTML updates <htmlFile>.html by parsing the template name <htmlFile>.gohtml.
func (s *server) updateHTML(htmlFile string) {
	t, err := template.ParseFiles(htmlFile)
	if err != nil {
		raiseError(err, ERROR)
	}
	fname := "./www/index.html"
	f, err := os.Create(fname)
	if err != nil {
		raiseError(err, ERROR)
	}
	if err := t.Execute(f, nil); err != nil {
		raiseError(err, ERROR)
	}
}

func (s *server) serve(url string) {
	r := mux.NewRouter()
	sub := r.Host("www.ormaybemore.com").Subrouter()
	sub.Handle("/", http.FileServer(http.Dir("./www")))

}

func main() {
	errlog = make(chan error, 3)
	errlogger := new(errorLogger)

	// start processing loop
	go errlogger.start(os.Stderr)

	Server.updateHTML("index.gohtml")
	Server.serve(":80")
}
