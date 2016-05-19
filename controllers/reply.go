package controllers

import (
	"github.com/astaxie/beego"
	"blog/models"
)

type ReplyControllers struct  {
	beego.Controller
}

func (c *ReplyControllers) Get(){

}
func (c *ReplyControllers) Post(){

}
func (c *ReplyControllers) Delete(){
	if !checkAcc(c.Ctx){
		c.Redirect("/login",302)
		return
	}
	tid := c.Input().Get("tid")
	err := models.DeleteReply(c.Input().Get("rid"),tid)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic/view?tid="+tid,302)
}
func (c *ReplyControllers) Add(){
	tid := c.Input().Get("tid")
	name := c.Input().Get("guest_name")
	content := c.Input().Get("content")
	err := models.AddReply(tid,name,content)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic/view?tid="+tid,302)
}