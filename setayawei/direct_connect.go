package page

import (
	"context"
	"net/http"
	"net/url"
	"time"

	uuid "github.com/satori/go.uuid"
)

type DirectConnectConfig []struct{ setting, value string }

// NewDirectConnectSession creates a direct video chat session, returning its
// joinable url, the key, and any errors.
func (s *Server) NewDirectConnectSession(ctx context.Context /* conf DirectConnectConfig */) (joinURL *url.URL, joinKey []byte, err error) {
	// TODO
	if d, ok := ctx.Deadline(); ok {
		if time.Now().UnixNano() > d.UnixNano() {
			// deadline is expired
			return nil, nil, context.DeadlineExceeded
		}

	}

}

func QuikFaceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// TODO: Should we be using TLS?
	case "GET":
		// TODO
		// TODO: set a cookie and respond with a redirect request using
		// POST

	case "POST":
		if err := r.ParseForm(); err != nil {
			// bad request
			http.Error(w, "bad request: cant parse request data", http.StatusBadRequest)
			return
		}
		// fetch the session id from the url
		if seshID := r.Form.Get("qfsession"); seshID != "" {
			// get it as a uuid
			uid, err := uuid.FromString(seshID)
			if err != nil {

				http.Error(w, "bad request: cant parse request data", http.StatusBadRequest)
				return
			}
			// fetch the QfSession's index in the pkg local map
			if sidx, ok := activeServer.QfSessions[uid]; ok {
				// a QfSession with this uuid has been found
				// so fetch it.
				sid := activeServer.qfsessions[sidx]

				// and get the QfSession from the QfAPI
				// interface.
				sesh := sid.Session()
				// is the session still valid (not expired), and
				// room for another participant?
				if !sesh.Valid(time.Now()) {
					// session expired

					// BUG: not sure if this is the correct
					// http status code...
					http.Error(w, "session expired", http.StatusServiceUnavailable)
					return
				}
				if sesh.UserCount < sesh.UserLimit {

					// TODO: check cookie auth before
					// passing to the handler.

				}
				// server full
				http.Error(w, "session user limit has been reached!", http.StatusForbidden)
				return

			}
		}
		// bad session id
		http.Error(w, "session not found", http.StatusNotFound)
		return

	// some other http method
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}

}
