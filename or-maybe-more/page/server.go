// +import "github.com/kdevb0x/or-maybe-more/src/page"
package page

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

var (
	// root dir to serve static assets from (including html)
	StaticAssetDir string

	// paths needed for TLS.
	// for dev only
	// *NOT FOR PRODUCTION*
	devcert string
	devkey  string
)

var DefaultServer = new(Server)

type Server struct {
	http.Server
	Active bool
}

// updateHTML updates <htmlFile>.html by parsing the template name <htmlFile>.gohtml.
func (s *Server) UpdateHTML(htmlfile string) error {
	info, err := os.Stat(htmlfile)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("cannot parse %s, file doesn't seem to exist!", htmlfile))
		}
	}
	t, err := template.ParseFiles(htmlfile)
	if err != nil {
		throw(err, ERROR)
	}
	if wwwdir := os.Getenv("OMM_WWW_DIR"); wwwdir != "" && StaticAssetDir == "" {
		StaticAssetDir = wwwdir
	}

	fname := filepath.Join(StaticAssetDir, info.Name()+".html")
	f, err := os.Create(fname)
	if err != nil {
		// throw(err, ERROR)
		return err
	}
	if err := t.Execute(f, nil); err != nil {
		// throw(err, ERROR)
		return err
	}
	return nil
}

func ContactFormHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		throw(err, ERROR)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func NotificationFormHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		throw(err, ERROR)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	e := r.Form.Get("email")
	if e != "" {
		err := addEmail(r.Context(), e, ContactInfo{})
		if err != nil {
			throw(err, ERROR)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func (s *Server) Serve(addr string) {
	r := mux.NewRouter()
	sub := r.Host("www.ormaybemore.com").Subrouter()
	sub.Handle("/", http.FileServer(http.Dir(filepath.Join(StaticAssetDir, "www"))))
	s.Handler = r
	s.Active = true
	log.Fatal(s.ListenAndServeTLS(devcert, devkey))
}

type NotificationListRow struct {
	EmailAddr string    `json:"email_addr"`
	Joined    time.Time `json:"joined,omitempty"`
}

func dbCreateNotificationListTable() error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	tx, err := conn.Begin(context.TODO())
	if err != nil {
		return err
	}
	defer tx.Rollback(nil)

	_, err = tx.Prepare(context.TODO(), "make_notifaction_list",
		`CREATE TABLE notification_list(
		email_addr VARCHAR,
		joined DATE
	) IF NOT EXIST;`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), "make_notifaction_list")
	if err != nil {
		return err
	}
	return nil
}

func addEmail(ctx context.Context, email string) error {
	var cancelCtx, cancelFn = context.WithCancel(ctx)
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			cancelFn()
		case <-done:
			return
		}
	}()
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	tx, err := conn.Begin(cancelCtx)
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(ctx, "add_email_addr",
		// TODO
		``)
	if err != nil {
		return err
	}
	defer conn.Deallocate(context.Background(), "add_email_addr")
	cmd, err := tx.Query(cancelCtx, stmt.SQL, email)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	if cmd.Err() != nil {
		return cmd.Err()
	}
	cmd.Close()
	done <- struct{}{}
	return nil

}
