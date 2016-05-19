package main

import (
	_ "blog/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"blog/models"
	"blog/controllers"
)
func init(){
	models.RegisterDB()
}
func main() {
	orm.Debug = true
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login",&controllers.LoginControllers{})
	beego.Router("/category",&controllers.CategoryController{})
	beego.Router("/topic",&controllers.TopicController{})
	beego.Router("/reply",&controllers.ReplyControllers{})
	beego.Router("/reply/add", &controllers.ReplyControllers{}, "post:Add")
	beego.Router("/reply/delete", &controllers.ReplyControllers{}, "get:Delete")
	beego.AutoRouter(&controllers.TopicController{})
	beego.Router("/topic/view",&controllers.TopicViewController{})
	orm.RunSyncdb("default",false,true)
	beego.Run()
}

