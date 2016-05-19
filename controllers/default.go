package controllers

import (
	"github.com/astaxie/beego"
	"blog/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	var err error
	c.Data["IsHome"] = true
	c.TplName = "home.html"
	c.Data["IsLogin"] = checkAcc(c.Ctx)
	c.Data["Categories"],err = models.GetCategory()
	if err != nil {
		beego.Error(err)
	}


	//分类显示
	show := c.Input().Get("show")
	label := c.Input().Get("label")
	c.Data["topics"],err = models.GetTopic(true,show,label)

}
