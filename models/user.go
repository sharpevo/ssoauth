package models

import (
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"ssoauth/mongo"
	"strings"
)

type User struct {
	Id       bson.ObjectId     `json:"id" bson:"_id"`
	Name     string            `json:"name"`
	Email    string            `json:"Email"`
	Password string            `json:"password`
	Roles    map[string]bool   `json:"roles"`
	Apps     map[string]bool   `json:"apps"`
	Tokens   map[string]string `json:"-"`
}

func (user *User) FindById(id string) (code int, err error) {
	session, err := mongo.CopyMasterSession()
	if err != nil {
		return ERROR_DATABASE, err
	}

	collection := session.DB(mongo.MongoConfig.Database).C("user")

	if !bson.IsObjectIdHex(id) {
		return ERROR_INPUT, err
	}

	err = collection.FindId(bson.ObjectIdHex(id)).One(user)
	if err != nil {
		return ERROR_NOT_FOUND, err
	}

	return 0, nil
}

func (user *User) HasToken(iss string, token string) bool {
	return user.Tokens[iss] == token
}

func (user *User) HasDomain2(domain string) bool {
	if user.Email == "quwubin@gmail.com" ||
		domain == beego.AppConfig.DefaultString("uicdomain", "accounts.igenetech.com") {
		return true
	} else {
		return user.Apps[strings.Replace(domain, ".", " ", -1)]
	}
}

func (user *User) HasDomain(domain string) bool {
	if user.Email == "quwubin@gmail.com" ||
		domain == beego.AppConfig.DefaultString(
			"uicdomain",
			"accounts.igenetech.com") {
		beego.Debug("CheckDomain: true, admin or ask for uic.")
		return true
	} else {
		app := uic.App{}

		session, err := mongo.CopyMasterSession()
		if err != nil {
			return false
		}
		collection := session.DB(mongo.MongoConfig.Database).C("app")
		err = collection.Find(bson.M{"domain": domain}).One(&app)
		beego.Debug("CheckDomain:", user.Apps[app.Id.Hex()], domain, "in", user.Apps)
		if err != nil {
			return false
		}
		return user.Apps[app.Id.Hex()]
	}
}
