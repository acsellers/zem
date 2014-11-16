package main

import (
	"net/http"
	"path/filepath"

	"github.com/acsellers/platform/router"
)

type AssetCtrl struct {
	*router.BaseController
}

func (AssetCtrl) Dupe() router.Controller {
	return AssetCtrl{&router.BaseController{}}
}

func (AssetCtrl) Path() string {
	return "assets"
}

func (sc AssetCtrl) Show() error {
	_, n := filepath.Split(sc.Params[":assetsid"])
	fn := filepath.Join(Config.ResourcePath, "assets", n)
	http.ServeFile(sc.Out, sc.Request, fn)
	return nil
}
