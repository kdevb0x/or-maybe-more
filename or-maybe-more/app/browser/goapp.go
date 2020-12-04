// +build ignore 

package ui

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// loginPage is the client oauth2 login page.
type loginPage struct {
	authData *authdata
}

type authdata struct {
	authobj   objx.Map
	hasCookie bool
}

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// some other error
	} else {
		// successful auth
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func initOmniauth() {
	// gomniauth.SetSecurityKey(/* TODO: add base64 encoded crypto key */)
	if privateKey := os.Getenv("OMM_SERVER_PRIVATE_KEY"); privateKey != "" {

		gomniauth.SetSecurityKey(privateKey)
	}
	gomniauth.WithProviders(
		facebook.New("key", "secret",
			"http://localhost:8080/auth/callback/facebook"),
		github.New("key", "secret",
			"http://localhost:8080/auth/callback/github"),
		google.New("key", "secret",
			"http://localhost:8080/auth/callback/google"),
	)
}

func authReqHandler(w http.ResponseWriter, r *http.Request, authchan ...chan authdata) {
	// format auth/{action}/{provider}

	vars := mux.Vars(r)
	action := vars["action"]
	provider := vars["provider"]
	switch action {
	case "login":
		pvdr, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		loginURL, err := pvdr.GetBeginAuthURL(nil, nil)
		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		pvdr, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		creds, err := pvdr.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		user, err := pvdr.GetUser(creds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		authcookie := objx.New(map[string]interface{}{
			"name": user.Name(),
		})
		authobj := authdata{
			authobj: authcookie,
		}
		authCookieVal := authcookie.MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieVal,
			Path:  "/",
		})
		if len(authchan) > 0 {
			authchan[0] <- authobj
		}
		w.Header()["Location"] = []string{"/session"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	}

}
