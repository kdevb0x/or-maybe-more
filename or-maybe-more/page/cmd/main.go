package main

import (
	"os"

	"../../page"
)

func main() {
	s := page.DefaultServer

	errlog = make(chan page.TaggedErr, 3)
	errlogger := new(page.ErrorLogger)

	// start processing loop
	go errlogger.Run(os.Stderr)

	s.UpdateHTML("index.gohtml")
	s.Serve(":8080")
}
