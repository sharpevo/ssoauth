package main

import (
	"github.com/astaxie/beego"
	"ssoauth/mongo"
	_ "ssoauth/routers"
)

func main() {
	err := mongo.Startup()
	if err != nil {
		beego.Error(err)
	}

	beego.Run()

	err = mongo.Shutdown()
	if err != nil {
		beego.Error(err)
	}
}
