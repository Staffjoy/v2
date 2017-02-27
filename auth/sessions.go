package auth

import (
	"net/http"
	"time"

	"v2.staffjoy.com/crypto"
)

// LoginUser sets a cookie to log in a user
func LoginUser(uuid string, support, rememberMe bool, res http.ResponseWriter) {
	var dur time.Duration
	var maxAge int

	if rememberMe {
		// "Remember me"
		dur = longSession
		maxAge = 0
	} else {
		dur = shortSession
		maxAge = int(dur.Seconds())
	}
	token, err := crypto.SessionToken(uuid, signingSecret, support, dur)
	if err != nil {
		panic(err)
	}
	cookie := &http.Cookie{
		Name:   cookieName,
		Value:  token,
		Path:   "/",
		Domain: "." + config.ExternalApex,
		MaxAge: maxAge,
	}
	http.SetCookie(res, cookie)
}

func getSession(req *http.Request) (uuid string, support bool, err error) {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		return
	}
	uuid, support, err = crypto.RetrieveSessionInformation(cookie.Value, signingSecret)
	return
}

// Logout forces an immediate logout of the current user.
// For internal applications - typically do an HTTP redirect
// to the www service `/logout/` route to trigger a logout.
func Logout(res http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
		Domain: "." + config.ExternalApex,
	}
	http.SetCookie(res, cookie)
}
