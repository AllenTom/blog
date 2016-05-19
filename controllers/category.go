package controllers

import (
	"github.com/astaxie/beego"
	"blog/models"
)

type CategoryController struct  {
	beego.Controller
}
func (c *CategoryController) Get(){
	op := c.Input().Get("op")
	switch op {
	case "add":
		name := c.Input().Get("name")
		if len(name) == 0{
			break
		}
		beego.Informational("name = ",name)
		err := models.AddCategory(name)
		if err != nil {
			beego.Error(err)
		}
		c.Redirect("/category",302)
		return
	case"del":
		id := c.Input().Get("id")
		if len(id) == 0{
			break
		}
		err := models.DeleteCategoty(id)
		if err != nil {
			beego.Error(err)
		}
		c.Redirect("/category",302)
		return

	}
	c.TplName = "category.html"
	var err error
	c.Data["Categories"] ,err = models.GetCategory()
	if err != nil {
		beego.Error(err)
	}
	c.Data["IsLogin"] = checkAcc(c.Ctx)
	c.Data["IsCategory"] = true

}

