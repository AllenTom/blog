package controllers

import (
	"github.com/astaxie/beego"
	"blog/models"
	"strings"
)

type TopicViewController struct {
	beego.Controller
}
func (c *TopicViewController) Get(){
	c.TplName = "test.html"
	c.Data["IsLogin"] = checkAcc(c.Ctx)
	tid := c.Input().Get("tid")
	var err error
	topic,err := models.GetTopicView(tid)
	c.Data["Topic"] = topic
	c.Data["Labels"] = strings.Split(topic.Labels," ")
	c.Data["Id"] = topic.Id
	c.Data["Replys"],err = models.GetReplys(topic.Id)
	if err != nil {
		beego.Error(err)
	}

}
