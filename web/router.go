package main

import (
	"github.com/acsellers/platform/router"
)

func Router() *router.Router {
	r := router.NewRouter()
	r.Many(AssetCtrl{})

	admin := r.Namespace("admin")
	admin.One(AdminSessionCtrl{})
	admin.Many(AdminBlogCtrl{})

	return r
}
