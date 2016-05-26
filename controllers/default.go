package controllers

import (
	"github.com/astaxie/beego"
	"blog/models"
	"strconv"
	"blog/controllers/ExtraUnit"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	controllers.InitTimeFunc()
	var err error
	c.Data["IsHome"] = true
	c.TplName = "home.html"
	c.Data["IsLogin"] = checkAcc(c.Ctx)
	topcount,err := models.GetTopic(true,"","")
	c.Data["AllCount"] = len(topcount)
	c.Data["Categories"],err = models.GetCategory()
	if err != nil {
		beego.Error(err)
	}
	//获取账户信息
	c.Data["Account"] = models.ReadUserInfo()

	//分类显示
	show := c.Input().Get("show")
	label := c.Input().Get("label")
	topics,err := models.GetTopic(true,show,label)
	c.Data["topics"] = topics

}

