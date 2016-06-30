package routers

import (
	"github.com/astaxie/beego"
	"ssoauth/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	nsCookie := beego.NewNamespace(
		"v1",
		beego.NSRouter(
			"/auth",
			&controllers.AuthByCookieController{}),
	)
	ns := beego.NewNamespace(
		"v2",
		beego.NSRouter(
			"/auth",
			&controllers.AuthController{}),
	)
	beego.AddNamespace(nsCookie)
	beego.AddNamespace(ns)
}
