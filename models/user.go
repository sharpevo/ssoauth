package models

import (
	"gopkg.in/mgo.v2/bson"
	"ssoauth/mongo"
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

func (user *User) HasDomain(domain string) bool {
	return user.Apps[strings.Replace(domain, ".", " ", -1)]
}
