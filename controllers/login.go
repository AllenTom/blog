package controllers

import (
"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type LoginControllers struct  {
	beego.Controller
}
func (c *LoginControllers) Get(){
	isExit := c.Input().Get("exit") == "true"
	if isExit{
		c.Ctx.SetCookie("username","",-1,"/")
		c.Ctx.SetCookie("pwd","",-1,"/")
	}
	c.TplName = "login.html"

}
func (c *LoginControllers) Post(){
	username := c.Input().Get("username")
	pwd := c.Input().Get("pwd")
	autoLogin := c.Input().Get("autologin") == "remember-me"

	if beego.AppConfig.String("username") == username &&
		beego.AppConfig.String("pwd") == pwd{
		maxAge := 0
		if autoLogin{
			maxAge = 1<<31 -1
		}
		c.Ctx.SetCookie("username",username,maxAge,"/")
		c.Ctx.SetCookie("pwd",pwd,maxAge,"/")
	}
	c.Redirect("/",301)
	return
}
func checkAcc(ctx *context.Context) bool{
	ck, err := ctx.Request.Cookie("username")
	if err != nil{
		return false
	}
	username := ck.Value

	ck, err = ctx.Request.Cookie("pwd")
	if err != nil{
		return false
	}
	pwd := ck.Value

	return (username == beego.AppConfig.String("username")  && pwd ==  beego.AppConfig.String("pwd"))
}