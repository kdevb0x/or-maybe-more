package main

import (
	"crypto/tls"
	"html/template"
	"net/http"
)

type server struct {
	cert tls.Certificate
	*http.Server
	// the root url path
	rootURL string
	// directory tree for the site
	pages map[string]template.Template // TODO: maybe use template.Template values?
}

func NewServer(addr string, certPath string) *server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", ServeTempPage)
	return &server{mux, listenaddr, make(map[string][]byte)}
}

func (s *server) AddTemplate(path string, t template.Template) {
	s.pages[path] = t
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// TODO implement caching
		if page, ok := s.pages[r.URL.String()]; ok {
			page.Execute(w, nil)
		}
	}
}
func (s *server) LoadCert(cert, key string) error {
	c, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return err
	}
	s.cert = c
	return nil

}

func ServeTempPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("check back soon!"))
}

var certkey string

func main() {
	s := NewServer("ormaybemore.com")
	s.ServeTLS()

}
