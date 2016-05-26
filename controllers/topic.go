package controllers

import (
	"github.com/astaxie/beego"
	"blog/models"
	"path"
	"blog/controllers/ExtraUnit"
)

type TopicController struct  {
	beego.Controller

}

func (c *TopicController) Get(){
	controllers.InitTimeFunc()
	c.Data["IsTopic"] = true
	c.TplName = "topic.html"
	c.Data["IsLogin"] = checkAcc(c.Ctx)
	var err error

	c.Data["Categories"],err = models.GetCategory()
	if err != nil {
		beego.Error(err)
	}
	cate := c.Input().Get("CategoryIndex")
	c.Data["Topics"] ,err = models.GetTopic(false,cate,"")
	if err != nil {
		beego.Error(err)
	}
	//获取账户信息
	c.Data["Account"] = models.ReadUserInfo()

}
func (c *TopicController) Post() {
	if !checkAcc(c.Ctx){
		c.Redirect("/login",302)
		return
	}
	title := c.Input().Get("title")
	content := c.Input().Get("content")
	category := c.Input().Get("category")
	label := c.Input().Get("label")
	//附件
	_,fh,err := c.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}
	var attachment string
	if fh != nil{
		attachment = fh.Filename
		err = c.SaveToFile("attachment",path.Join("attachment",attachment))
		if err != nil {
			beego.Error(err)
		}
	}

	tid := c.Input().Get("tid")
	if len(tid) != 0{
		models.ModifyTopic(tid,title,content,category,label,attachment)
	}else {
		err = models.AddTopic(title, content,category,label,attachment)
		if err != nil {
			beego.Error(err)
		}
	}
	c.Redirect("/topic",302)

}
func (c *TopicController) Add() {
	if !checkAcc(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	c.TplName = "add_topic.html"
}
func (c *TopicController) Modify() {
	if !checkAcc(c.Ctx){
		c.Redirect("/login",302)
	}
	c.Data["IsLogin"] = checkAcc(c.Ctx)
	c.TplName = "topic_modify.html"
	tid := c.Input().Get("tid")
	topic, err := models.GetTopicView(tid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
	}
	c.Data["Topic"] = topic
	c.Data["tid"] = tid
}
func (c *TopicController) Delete(){
	if !checkAcc(c.Ctx){
		c.Redirect("/login",302)
	}
	tid := c.Input().Get("tid")
	err := models.DeleteTopic(tid)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/",302)
}