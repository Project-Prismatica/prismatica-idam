package prismatica_idam


import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

const (
	sessionsName = "prismatica-auth-session"
)

var (
	// TODO: dynamically generate this
	sessionStore = sessions.NewCookieStore([]byte("initial-secret-key"))
)

func nullAuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"path": r.URL, "client": r.RequestURI}).
		Info("handling with null authentication")

	session, err := sessionStore.Get(r, sessionsName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
	log.WithFields(log.Fields{"session": session}).
		Info("client authorized")

}

func RunAmbassadorAuthenticationService(bindSpec string) (err error) {

	router := mux.NewRouter()

	httpServer := &http.Server{
		Handler: 		router,
		Addr:			bindSpec,
		WriteTimeout: 	15 * time.Second,
		ReadTimeout: 	15 * time.Second,
	}


	router.MatcherFunc(func(r *http.Request, m *mux.RouteMatch) bool {return true}).
		HandlerFunc(nullAuthenticationHandler)

	log.WithFields(log.Fields{"bind": bindSpec}).
		Debug("starting http server")
	serverRunError := httpServer.ListenAndServe()
	if serverRunError != nil {
		log.WithFields(log.Fields{"error": serverRunError}).
			Error("could not run server")
	}

	httpServer.Close()

	return
}
