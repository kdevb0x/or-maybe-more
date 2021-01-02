package app

import (
	"fmt"
	"net/http"
	"os"
	"sort"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	_ "github.com/markbates/goth/providers/instagram"
	_ "github.com/markbates/goth/providers/twitter"
)

var _ = gothic.SessionName

var ServerRootURL string

func initGoth() {
	goth.UseProviders(google.New(os.Getenv("OMM_GOOGLE_OAUTH_KEY"), os.Getenv("OMM_GOOGLE_OAUTH_SECRET"), fmt.Sprintf("%s/auth/google/callback", ServerRootURL)),
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://localhost:3000/auth/facebook/callback"))

	// instagram.New(os.Getenv("INSTAGRAM_KEY"), os.Getenv("INSTAGRAM_SECRET"), "http://localhost:3000/auth/instagram/callback"),
	// apple.New(os.Getenv("APPLE_KEY"), os.Getenv("APPLE_SECRET"), "http://localhost:3000/auth/apple/callback", nil, apple.ScopeName, apple.ScopeEmail),
	m := make(map[string]string)
	m["google"] = "Google"
	m["facebook"] = "Facebook"

	// m["instagram"] = "Instagram"
	// m["apple"] = "Apple"

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	providerIDX

}

// oauth data from provider
type OAuthUserData struct {
	logoutURL string
	Provider  string
	// FirstName and LastName
	Name         [2]string
	Email        string
	NickName     string
	Location     string
	AvatarURL    string
	Description  string
	UserID       string
	AccessToken  string
	ExpiresAt    string
	RefreshToken string
}

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// some other error
	} else {
		if token := c.Value; token != "" {
			println(fmt.Printf("auth cookie token : %s\n", token))
		}
		// previous successful auth
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

/*
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
*/
