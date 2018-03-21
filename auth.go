package prismatica_idam


import (
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

const (
	authCookieName = "prismatica-auth-session"
)

func generateAuthCookie()(cookie *http.Cookie) {
	cookie = &http.Cookie{Name: authCookieName,
		Value: "null-authentication"}
	return
}

func nullAuthenticationHandler(c echo.Context)(err error) {
	log.WithFields(log.Fields{"path": c.Path(), "client": c.RealIP()}).
		Info("handling with null authentication")

	authCookie, authCookieFetchError := c.Request().Cookie(authCookieName)
	if authCookieFetchError != nil {
		log.WithFields(log.Fields{"authCookieFetchError": authCookieFetchError,
			}).Debug("did not find cookie")
		authCookie = generateAuthCookie()
		c.SetCookie(authCookie)
	}

	err = c.String(http.StatusOK, "")

	log.WithFields(log.Fields{"auth_cookie": authCookie}).
		Info("client authorized")

	return
}

func RunAmbassadorAuthenticationService(bindSpec string) (err error) {

	httpServer := echo.New()

	httpServer.Any("/", nullAuthenticationHandler)

	log.WithFields(log.Fields{"bind": bindSpec}).
		Debug("starting http server")
	serverRunError := httpServer.Start(bindSpec)
	if serverRunError != nil {
		log.WithFields(log.Fields{"error": serverRunError}).
			Error("could not run server")
	}

	return
}
