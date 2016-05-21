package controllers

import (
	"github.com/astaxie/beego"
)

type AccountControllers struct  {
	beego.Controller
}

func (c *AccountControllers) Get() {
	c.TplName = "account_setting.html"
	c.Data["IsLogin"] = checkAcc(c.Ctx)
}
