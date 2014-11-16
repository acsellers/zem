package main

import (
	"net/http"

	"github.com/acsellers/platform/router"
)

type SessionCtrl struct {
	*BaseCtrl
}

func (sc SessionCtrl) Dupe() router.Controller {
	return SessionCtrl{DupeBaseCtrl()}
}

func (SessionCtrl) Path() string {
	return "sessions"
}

func (sc SessionCtrl) Show() error {
	if sc.Current.User != nil {
		return router.RedirectError{
			Location: "/",
			Reason:   "Already logged in",
		}
	}

	sc.Current.Template = "sessions_show.html"
	return sc.Render()
}

func (sc SessionCtrl) Update() error {
	if sc.Current.User != nil {
		return router.RedirectError{
			Location: "/",
			Reason:   "Already logged in",
		}
	}
	sc.Request.ParseForm()
	// Change if Conn is becoming special
	sc.Current.User = Conn.User.Authenticate(sc.Request.Form.Get("username"), sc.Request.Form.Get("password"))
	if sc.Current.User == nil {
		sc.Log.Println("Could not retrieve user")
		return router.RedirectError{
			Location: "/sessions",
			Reason:   "Error Finding User",
		}
	}
	enc, err := Cookie.Encode("zem-login", sc.Current.User.ID)
	if err != nil {
		return router.RedirectError{
			Location: "/sessions",
			Reason:   "Error encoding cookie: " + err.Error(),
		}
	}

	sc.Log.Printf("Setting cookie for user (%v)\n", sc.Current.User.ID)
	http.SetCookie(sc.Out, &http.Cookie{
		Name:  "zem-login",
		Value: enc,
		Path:  "/",
	})

	return router.RedirectError{
		Location: sc.Request.Referer(),
		Reason:   "Redirecting to previous",
	}
}

func (sc SessionCtrl) Delete() error {
	http.SetCookie(sc.Out, &http.Cookie{Name: "zem-login", Value: "", Path: "/"})
	return router.RedirectError{
		Location: "/",
		Reason:   "Logged out User",
	}
}
