package main

import (
	"github.com/acsellers/platform/router"
	"github.com/acsellers/zem/store"
)

type BaseCtrl struct {
	*router.BaseController
	Conn    *store.Conn
	Err     error
	Current struct {
		User     *store.User
		Template string
	}
}

func DupeBaseCtrl() *BaseCtrl {
	return &BaseCtrl{
		BaseController: &router.BaseController{},
		Conn:           Conn.Clone(),
	}
}

func (bc *BaseCtrl) PreFilter() error {
	bc.Conn.Log = bc.Log
	bc.Context = make(map[string]interface{})
	var value string
	if cookie, err := bc.Request.Cookie("zem-login"); err != nil {
		bc.Log.Printf("Could not retrieve Cookie: %s\n", err.Error())
		return nil
	} else {
		value = cookie.Value
	}

	var userid int
	if err := Cookie.Decode("zem-login", value, &userid); err != nil {
		bc.Log.Printf("Could not decode Cookie: %s\n", err.Error())
	}

	var u store.User
	u, bc.Err = bc.Conn.User.Eq(userid).Retrieve()
	if bc.Err != nil {
		bc.Log.Println("Could not retrieve user for cookie")
		return nil
	}
	bc.Current.User = &u
	bc.Context["CurrentUser"] = u

	return nil
}

func (bc *BaseCtrl) Render() error {
	bc.Out.Header().Add("Content-Type", "text/html")
	bc.Err = Template.ExecuteTemplate(bc.Out, bc.Current.Template, bc.Context)
	if bc.Err != nil {
		bc.Log.Println("Encountered error during template execution:", bc.Err)
		return bc.Err
	}
	return nil
}

func (bc *BaseCtrl) FormString(key string) string {
	return bc.Request.Form.Get(key)
}
