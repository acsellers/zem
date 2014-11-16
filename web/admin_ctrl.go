package main

import (
	"net/http"

	"github.com/acsellers/platform/router"
)

type AdminSessionCtrl struct {
	SessionCtrl
}

func (asc AdminSessionCtrl) Dupe() router.Controller {
	return AdminSessionCtrl{SessionCtrl{DupeBaseCtrl()}}
}

func (asc AdminSessionCtrl) Show() error {
	if asc.Current.User != nil {
		return router.RedirectError{
			Location: "/admin/blogs",
			Reason:   "Already logged in",
		}
	}
	return asc.SessionCtrl.Show()
}

func (asc AdminSessionCtrl) Update() error {
	if asc.Current.User != nil {
		return router.RedirectError{
			Location: "/admin/blogs",
			Reason:   "Already logged in",
		}

	}
	asc.Request.ParseForm()
	// Change if Conn is becoming special
	asc.Current.User = Conn.User.
		Management().Eq(true).
		Authenticate(asc.Request.Form.Get("username"), asc.Request.Form.Get("password"))
	if asc.Current.User == nil {
		asc.Log.Println("Could not retrieve user")
		return router.RedirectError{
			Location: "/sessions",
			Reason:   "Error Finding User",
		}
	}
	enc, err := Cookie.Encode("zem-login", asc.Current.User.ID)
	if err != nil {
		return router.RedirectError{
			Location: "/sessions",
			Reason:   "Error encoding cookie: " + err.Error(),
		}
	}

	asc.Log.Printf("Setting cookie for user (%v)\n", asc.Current.User.ID)
	http.SetCookie(asc.Out, &http.Cookie{
		Name:  "zem-login",
		Value: enc,
		Path:  "/",
	})

	return router.RedirectError{
		Location: "/admin/blogs",
		Reason:   "Redirecting to previous",
	}
}

type AdminBlogCtrl struct {
	*BaseCtrl
}

func (abc AdminBlogCtrl) Dupe() router.Controller {
	return AdminBlogCtrl{DupeBaseCtrl()}
}

func (abc AdminBlogCtrl) Path() string {
	return "blogs"
}

func (abc AdminBlogCtrl) Index() error {
	blargs, err := abc.Conn.Blog.RetrieveAll()
	if err != nil {
		return err
	}
	abc.Context["Blogs"] = blargs
	abc.Current.Template = "admin_blogs_index.html"
	return abc.Render()
}

func (abc AdminBlogCtrl) Show() error {
	return nil
}

func (abc AdminBlogCtrl) New() error {
	return nil
}

func (abc AdminBlogCtrl) Create() error {
	return nil
}

func (abc AdminBlogCtrl) Edit() error {
	return nil
}

func (abc AdminBlogCtrl) Update() error {
	return nil
}

func (abc AdminBlogCtrl) Delete() error {
	return nil
}
