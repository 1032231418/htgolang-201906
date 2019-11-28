package routers

import (
	"github.com/astaxie/beego"
	"github.com/xxdu521/cmdbgo/server/controllers"
	"github.com/xxdu521/cmdbgo/server/controllers/auth"
)

func init(){
	beego.AutoRouter(&auth.AuthController{})
	//beego.AutoRouter(&controllers.TestController{})
	beego.AutoRouter(&controllers.UserPageController{})
	beego.AutoRouter(&controllers.UserController{})
	beego.AutoRouter(&controllers.TokenController{})
}
