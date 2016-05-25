package models

import (
	"os"
	"time"
	"path"
	"strings"
	_"github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	_"github.com/mattn/go-sqlite3"
	"strconv"
	"github.com/astaxie/beego"
	"encoding/json"
)

const (
	_MYSQL = "mysql"
)

type Account struct {
	UserAvatar    string
	Name          string
	Sex           int64
	BirthdayYear  int64
	BirthdayMonth int64
	BirthdayDay   int64
	Organization  string
	Position      string
	Qq_number     string
	Email         string
	Google_id     string
	Google_web    string
	Twitter_id    string
	Twitter_web   string
	Facebook_id   string
	Facebook_web  string
	Github_id     string
	Github_web    string
	Weibo_id      string
	Weibo_web     string
	Linkedin_id   string
	Linkedin_web  string
	Tumblr_id     string
	Tumblr_web    string
	Steam_id      string
	Steam_web     string
	Home_message  string
	Private_intro string
}

type Category struct {
	Id              int64
	Title           string
	Created         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	TopicTime       time.Time `orm:"index"`
	TopicCount      int64
	TopicLastUserId int64
}
type Reply struct {
	Id      int64
	Tid     int64
	Name    string
	Content string`orm:"size(5000)"`
	Time    time.Time`orm:"index"`
}
type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Content         string`orm:"size(5000)"`
	Attachment      string
	Created         time.Time `orm:"index"`
	Updated         time.Time `orm:"index"`
	Category        string
	TopicType       string
	Labels          string
	Views           int64     `orm:"index"`
	Author          string
	ReplyTime       time.Time `orm:"index"`
	ReplyCount      int64
	ReplyLastUserId int64
}
type User struct {
	Id           int64
	UserName     string
	Pwd          string
	NickName     string
	Email        string
	PersonalInfo string
}
type Tag struct {
	Id         int64
	Name       string
	TopicCount int64
}

func AddReply(tid, name, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	reply := &Reply{
		Name:name,
		Tid:tidNum,
		Content:content,
		Time:time.Now(),
	}
	_, err = o.Insert(reply)
	if err != nil {
		return err
	}
	//取出对应用户id
	replys := new(Reply)
	qss := o.QueryTable("reply")
	err = qss.Filter("name", name).One(replys)
	LastReplyId := replys.Id
	topic := &Topic{
		Id : tidNum,
	}
	if o.Read(topic) == nil {
		topic.ReplyTime = time.Now()
		topic.ReplyCount++
		topic.ReplyLastUserId = LastReplyId
		_, err = o.Update(topic)
	}
	return err
}
func DeleteReply(rid, tid string) error {
	ridNum, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		return err
	}
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	reply := &Reply{Id:ridNum}
	_, err = o.Delete(reply)
	if err != nil {
		return err
	}
	topic := &Topic{
		Id : tidNum,
	}
	if o.Read(topic) == nil {
		topic.ReplyCount--
		_, err = o.Update(topic)
	}
	return err
}
func GetTopicView(tid string) (*Topic, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	topic := new(Topic)
	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidNum).One(topic)
	topic.Views++
	_, err = o.Update(topic)
	if err != nil {
		return nil, err
	}
	topic.Labels = strings.Replace(strings.Replace(strings.Replace(
		topic.Labels, "#$", " ", -1), "$", "", -1), "#", "", -1)
	return topic, err
}
func GetReplys(tid int64) ([]*Reply, error) {

	o := orm.NewOrm()
	replys := make([]*Reply, 0)
	qs := o.QueryTable("reply")
	_, err := qs.Filter("tid", tid).OrderBy("-time").All(&replys)
	return replys, err
}
func AddTopic(title, content, category, labels, attachment string) error {
	//tags
	labels = "$" + strings.Join(strings.Split(labels, " "), "#$") + "#"

	o := orm.NewOrm()
	topic := &Topic{
		Title:title,
		Content:content,
		Created:time.Now(),
		Updated:time.Now(),
		Attachment:attachment,
		Category:category,
		Labels:labels,
		ReplyTime:time.Now(),

	}
	_, err := o.Insert(topic)
	err = AddCategory(category)
	if err != nil {
		return err
	}
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		// 如果不存在我们就直接忽略，只当分类存在时进行更新
		cate.TopicCount++
		_, err = o.Update(cate)
	} else {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", category).One(cate)
		if err == nil {
			cate.TopicCount++
			_, err = o.Update(cate)
		}
	}
	return err
}
func ModifyTopic(tid, title, content, category, labels, attachment string) error {
	labels = "$" + strings.Join(strings.Split(labels, " "), "#$") + "#"
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	var elderCate, elderAttachment string
	o := orm.NewOrm()
	topic := &Topic{
		Id : tidNum,
	}
	if o.Read(topic) == nil {
		elderCate = topic.Category
		elderAttachment = topic.Attachment
		topic.Title = title
		topic.Content = content
		topic.Category = category
		topic.Attachment = attachment
		topic.Labels = labels
		topic.Updated = time.Now()
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
	}
	// 更新分类统计
	if len(elderCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", elderCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_, err = o.Update(cate)
		}
	}

	//删除附件
	if len(elderAttachment) > 0 {
		os.Remove(path.Join("attachment", elderAttachment))
	}

	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		cate.TopicCount++
		_, err = o.Update(cate)
	}
	return nil

	return err
}
func GetTopic(isDesc  bool, category, label string) ([]*Topic, error) {
	o := orm.NewOrm()
	topics := make([]*Topic, 0)
	qs := o.QueryTable("topic")
	var err error
	if isDesc {
		if len(category) > 0 {
			qs = qs.Filter("category", category)
		}
		if len(label) > 0 {
			qs = qs.Filter("labels__contains", "$" + label + "#")
		}
		_, err = qs.OrderBy("-created").All(&topics)
	} else {
		_, err = qs.All(&topics)
	}

	return topics, err
}
func RegisterDB() {
	var _DB_USERNAME = beego.AppConfig.String("mysqluser")
	var _DB_PWD = beego.AppConfig.String("mysqlpwd")
	var _DB_URL = beego.AppConfig.String("mysqlurl")
	var _DB_PORT = beego.AppConfig.String("mysqlport")
	var _DB_NAME = beego.AppConfig.String("mysqldb")
	orm.RegisterModel(new(Category), new(Topic), new(Reply), new(User), new(Tag))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", _DB_USERNAME + ":" + _DB_PWD + "@tcp(" + _DB_URL + ":" + _DB_PORT + ")/" + _DB_NAME + "?charset=utf8&loc=Asia%2FChongqing", 30, 200)
}
/*func SaveTopicTag(tid int64,tags string){
	taglist := strings.Split(tags," ")
	o := orm.NewOrm()
	o.QueryTable("tag")
	var err error
	for x := range taglist {
		tag := new(Tag)
		tag.Name = taglist[x]
		err = o.Read(&tag)
		if err == orm.ErrNoRows{
			o.Insert(tag)
		}else {
			tag.TopicCount++
			_,err = o.Update(tag)
		}
	}
}*/

func AddCategory(name string) error {
	o := orm.NewOrm()

	cate := &Category{
		Title:name,
		Created: time.Now(),
		TopicTime:time.Now(),
	}
	query := o.QueryTable("category")
	err := query.Filter("title", name).One(cate)
	if err == nil {
		return err
	}
	_, err = o.Insert(cate)
	if err != nil {
		return err
	}
	return nil
}
func GetCategory() ([]*Category, error) {
	o := orm.NewOrm()
	cates := make([]*Category, 0)
	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}
func DeleteCategoty(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := &Category{Id : cid}
	_, err = o.Delete(cate)
	return err

}
func DeleteTopic(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := &Topic{Id : cid}
	_, err = o.Delete(cate)

	return err
}
func SaveUserInfo(account Account) {
	userFile := "AccInfo.json"
	fout, err := os.Create(userFile)
	if err != nil {
		beego.Info(err)
	}
	b, err := json.Marshal(account)
	fout.Write(b)
}
func ReadUserInfo() *Account{
	jsonfile, err := os.Open("AccInfo.json")
	if err != nil {
		beego.Error(err)
	}
	Account := new(Account)
	UserInfo := json.NewDecoder(jsonfile)
	UserInfo.Decode(&Account)
	return Account
}