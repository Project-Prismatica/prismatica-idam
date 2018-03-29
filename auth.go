package prismatica_idam


import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	HeaderName = "x-prismatica-session"
	externalPathVarName = "external-path"
)

var (
	// TODO: get this elsewhere
	jwtSigningKey = []byte("soopersecret")
	jwtSigningMethod = jwt.SigningMethodHS256
)

func generateJwt()(token *jwt.Token) {
	token = jwt.New(jwtSigningMethod)
	return
}

func validateJwt(toValidate *jwt.Token)(interface{}, error) {
	// always provide the same signing key
	return jwtSigningKey, nil
}

func parseJwt(rawToken string)(token *jwt.Token, err error) {
	token, err = jwt.Parse(rawToken, validateJwt)
	return
}

func jwtAuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)
	externalPath := requestVars[externalPathVarName]
	requestId := r.Header.Get("x-request-id")
	currentTokenRaw := r.Header.Get(HeaderName)
	var currentToken *jwt.Token
	var tokenParseError error

	log.WithFields(log.Fields{"path": r.URL, "client": r.RequestURI,
			"externalPath": externalPath, "requestId": requestId}).
		Info("handling authentication")

	if 0 == len(currentTokenRaw) {
		log.WithFields(log.Fields{}).Info("creating new session")
		currentToken = generateJwt()
	} else {
		currentToken, tokenParseError = parseJwt(currentTokenRaw)
	}

	if tokenParseError != nil {
		log.WithFields(log.Fields{"jwtParseError": tokenParseError,
				"requestId": requestId}).
			Warn("could not parse JWT, assigning new")
		currentToken = generateJwt()
	}

	signedKey, keySigningError := currentToken.SignedString(jwtSigningKey)
	if keySigningError != nil {
		log.WithFields(log.Fields{"error": keySigningError,
				"requestId": requestId}).
			Error("could not sign JWT")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set(HeaderName, signedKey)

	log.WithFields(log.Fields{"requestId": requestId}).
		Info("client authorized")

}

func RunAmbassadorAuthenticationService(bindSpec string) (err error) {

	router := mux.NewRouter().
		StrictSlash(false)

	router.HandleFunc(
		fmt.Sprintf("extauth/{%s:.*}", externalPathVarName),
		jwtAuthenticationHandler)

	loggingMiddleware := handlers.CombinedLoggingHandler(os.Stdout, router)
	memoryManagementMiddleware := context.ClearHandler(loggingMiddleware)

	httpServer := &http.Server{
		Handler: 		memoryManagementMiddleware,
		Addr:			bindSpec,
		WriteTimeout: 	15 * time.Second,
		ReadTimeout: 	15 * time.Second,
	}

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
