package models

import (
	"os"
	"time"
	"path"
	"strings"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/orm"
	_"github.com/mattn/go-sqlite3"
	"strconv"
)

const (
	_DB_NAME = "Database/blog.db"
	_SQLITE3_DRIVER = "sqlite3"
)

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
	Labels			string
	Views           int64     `orm:"index"`
	Author          string
	ReplyTime       time.Time `orm:"index"`
	ReplyCount      int64
	ReplyLastUserId int64
}
func AddReply(tid,name,content string) error{
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return  err
	}
	o := orm.NewOrm()
	reply := &Reply{
		Name:name,
		Tid:tidNum,
		Content:content,
		Time:time.Now(),
	}
	_,err = o.Insert(reply)
	if err != nil {
		return err
	}
	//取出对应用户id
	replys := new(Reply)
	qss := o.QueryTable("reply")
	err = qss.Filter("name",name).One(replys)
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
func DeleteReply(rid,tid string) error{
	ridNum, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		return  err
	}
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return  err
	}
	o := orm.NewOrm()
	reply := &Reply{Id:ridNum}
	_,err = o.Delete(reply)
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
		topic.Labels, "#$", " ", -1), "$", "", -1),"#","",-1)
	return topic, err
}
func GetReplys(tid int64) ([]*Reply,error){

	o := orm.NewOrm()
	replys  := make([]*Reply,0)
	qs := o.QueryTable("reply")
	_,err := qs.Filter("tid",tid).OrderBy("-time").All(&replys)
	return replys,err
}
func AddTopic(title, content, category,labels string) error {
	//tags
	labels = "$" + strings.Join(strings.Split(labels," "), "#$") + "#"

	o := orm.NewOrm()
	topic := &Topic{
		Title:title,
		Content:content,
		Created:time.Now(),
		Updated:time.Now(),
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
	}else{
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title",category).One(cate)
		if err == nil {
			cate.TopicCount++
			_, err = o.Update(cate)
		}
	}

	return err
}
func ModifyTopic(tid, title, content, category,labels string) error {
	labels = "$" + strings.Join(strings.Split(labels," "),"#$") + "#"
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	var elderCate string
	o := orm.NewOrm()
	topic := &Topic{
		Id : tidNum,
	}
	if o.Read(topic) == nil {
		elderCate = topic.Category
		topic.Title = title
		topic.Content = content
		topic.Category = category
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
func GetTopic(isDesc  bool,category,label string) ([]*Topic, error) {
	o := orm.NewOrm()
	topics := make([]*Topic, 0)
	qs := o.QueryTable("topic")
	var err error
	if isDesc {
		if len(category) > 0 {
			qs = qs.Filter("category", category)
		}
		if len(label) > 0{
			qs = qs.Filter("labels__contains","$"+label+"#")
		}
		_, err = qs.OrderBy("-created").All(&topics)

	} else {
		_, err = qs.All(&topics)
	}

	return topics, err
}
func RegisterDB() {
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}
	orm.RegisterModel(new(Category), new(Topic),new(Reply))
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)
}

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