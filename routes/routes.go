package routes

import (
	"github.com/astaxie/beego"
	"github.com/naokij/imgdrop/controllers"
)

func InitRoutes() {

	beego.ErrorController(&controllers.ErrorController{})

	beego.Router("/", &controllers.MainController{})
	beego.Router("/upload", &controllers.UploadController{})
	beego.Router("/user", &controllers.UserController{})
	beego.Router("/my_images", &controllers.ImagesController{})

	//登录
	authController := new(controllers.AuthController)
	beego.Router("/login", authController, "get:Login;post:DoLogin")
	beego.Router("/login/:returnurl(.+)", authController, "get:Login")
	beego.Router("/logout", authController, "get:Logout")

	//管理用户
	beego.Router("/users/new", &controllers.UserController{}, "get:New")
	beego.Router("/users/create", &controllers.UserController{}, "post:Create")
	beego.Router("/users/:id:int/edit", &controllers.UserController{}, "get:Edit")
	beego.Router("/users/:id:int/update", &controllers.UserController{}, "post:Update")
	beego.Router("/users/:id:int/del", &controllers.UserController{}, "delete:Del")
	beego.Router("/users", &controllers.UserController{}, "get:List")
}
