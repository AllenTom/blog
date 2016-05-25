package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"blog/models"
	"path"
	"path/filepath"
)

type AccountControllers struct {
	beego.Controller
}

func (c *AccountControllers) Get() {
	c.TplName = "account_setting.html"
	c.Data["IsLogin"] = checkAcc(c.Ctx)
	c.Data["Account"] = models.ReadUserInfo()
}
func (c *AccountControllers) Post() {
	name := c.Input().Get("username")
	sex := c.Input().Get("sex")
	birth_year := c.Input().Get("birthday_year")
	birth_month := c.Input().Get("birthday_month")
	birth_day := c.Input().Get("birthday_day")
	organization := c.Input().Get("organization")
	position := c.Input().Get("position")
	qq_number := c.Input().Get("qq_number")
	email := c.Input().Get("mail")
	google_id := c.Input().Get("google_id")
	google_web := c.Input().Get("google_web")
	twitter_id := c.Input().Get("twitter_id")
	twitter_web := c.Input().Get("twitter_web")
	facebook_id := c.Input().Get("facebook_id")
	facebook_web := c.Input().Get("facebook_web")
	github_id := c.Input().Get("github_id")
	github_web := c.Input().Get("github_web")
	weibo_id := c.Input().Get("weibo_id")
	weibo_web := c.Input().Get("weibo_web")
	linkedin_id := c.Input().Get("linkedin_id")
	linkedin_web := c.Input().Get("linkedin_web")
	tumblr_id := c.Input().Get("tumblr_id")
	tumblr_web := c.Input().Get("tumblr_web")
	steam_id := c.Input().Get("steam_id")
	steam_web := c.Input().Get("steam_web")
	home_message := c.Input().Get("home_message")
	private_intro := c.Input().Get("private_intro")

	//附件
	_, fh, err := c.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}
	var attachment string
	if fh != nil {
		extension := filepath.Ext(fh.Filename)
		fh.Filename = "user_avatar" + extension
		attachment = fh.Filename
		err = c.SaveToFile("attachment", path.Join("static", "img", attachment))
		if err != nil {
			beego.Error(err)
		}
	}
	account := models.Account{
		UserAvatar:attachment,
		Name:name,
		Sex:toInt(sex),
		BirthdayYear:toInt(birth_year),
		BirthdayMonth:toInt(birth_month),
		BirthdayDay:toInt(birth_day),
		Organization:organization,
		Position:position,
		Qq_number : qq_number,
		Email : email,
		Google_id : google_id,
		Google_web : google_web,
		Twitter_id : twitter_id,
		Twitter_web : twitter_web,
		Facebook_id : facebook_id,
		Facebook_web : facebook_web,
		Github_id : github_id,
		Github_web : github_web,
		Weibo_id : weibo_id,
		Weibo_web : weibo_web,
		Linkedin_id : linkedin_id,
		Linkedin_web : linkedin_web,
		Tumblr_id : tumblr_id,
		Tumblr_web : tumblr_web,
		Steam_id : steam_id,
		Steam_web : steam_web,
		Home_message : home_message,
		Private_intro:private_intro,
	}
	models.SaveUserInfo(account)
	c.Redirect("/", 302)

}
func toInt(str string) int64 {
	tidNum, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		beego.Info(err)
	}
	return tidNum
}