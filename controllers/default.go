package controllers

import (
	"github.com/astaxie/beego"
	"blog/models"
	"strconv"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
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

	//分页
	page := c.Input().Get("page")
	pageNum, err := strconv.ParseInt(page, 10, 64)
	pages := make([]int,5,5)
	pageCount := (len(topics) / 10) + 1
	if (len(page) == 0){
		pageNum = 0
	}
	pageNum = 0
	for x := 0; x < 5 ;x++{
		pages[x] = int(pageNum) + x
		if(pages[x] == pageCount){
			break
		}
	}
}
