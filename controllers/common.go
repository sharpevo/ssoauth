package controllers

import (
	"crypto/rsa"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"strings"
)

type ControllerError struct {
	Status  int    `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	ErrInput   = &ControllerError{400, 10001, "Invalid Inputs"}
	ErrExpired = &ControllerError{400, 10012, "Token is expired"}
)

type BaseController struct {
	beego.Controller
}

type NestPreparer interface {
	NestPrepare()
}

var PubKey *rsa.PublicKey

func (c *BaseController) Prepare() {
	pubBytes, err := ioutil.ReadFile("keys/ip.rsa.pub")
	if err != nil {
		beego.Error("ReadPublicBytes:", err)
		return
	}
	PubKey, err = jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	beego.Debug("ReadPublicKeys")

	if app, ok := c.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}

func (c *BaseController) ParseUserAgent() (userAgent string) {
	userAgent = c.Ctx.Input.UserAgent()
	if userAgent == "wechat" {
		userAgent = "wechat"
	} else {
		userAgent = "browser"
	}

	return userAgent
}

func (c *BaseController) ParseToken(source string) (t *jwt.Token, e *ControllerError) {
	var tokenString string
	switch source {
	case "cookie":
		authString := c.Ctx.Input.Cookie("token")
		beego.Debug("AuthString:", authString)
		tokenString = authString
	case "header":
		authString := c.Ctx.Input.Header("Authorization")
		beego.Debug("AuthString:", authString)
		kv := strings.Split(authString, " ")
		if len(kv) != 2 || kv[0] != "Bearer" {
			beego.Error("AuthString invalid:", authString)
			return nil, ErrInput
		}
		tokenString = kv[1]
	default:
		c.AuthFail()
		return
	}

	pubBytes, err := ioutil.ReadFile("keys/ip.rsa.pub")
	if err != nil {
		beego.Error("ReadPublicBytes:", err)
		return
	}

	PubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		beego.Error("ParseRSAPublicKey:", err)
		return
	}

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return PubKey, nil
		},
	)

	if err != nil {
		beego.Error("ParseToken:", err)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// not a token
				return nil, ErrInput
			} else if ve.Errors&
				(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// expired or active yet
				return nil, ErrExpired
			} else {
				return nil, ErrInput
			}
		} else {
			return nil, ErrInput
		}
	}

	if !token.Valid {
		beego.Error("TokenInvalid:", tokenString)
		return nil, ErrInput
	}
	beego.Debug("Token:", token)
	return token, nil
}

func (c *BaseController) ParseTokenFromCookie() (t *jwt.Token, e *ControllerError) {
	authString := c.Ctx.Input.Header("Authorization")
	beego.Debug("AuthString:", authString)

	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		beego.Error("AuthString invalid:", authString)
		return nil, ErrInput
	}
	tokenString := kv[1]

	pubBytes, err := ioutil.ReadFile("keys/ip.rsa.pub")
	if err != nil {
		beego.Error("ReadPublicBytes:", err)
		return
	}
	PubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	beego.Debug("ReadPublicKeys:", PubKey)

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return PubKey, nil
		},
	)

	if err != nil {
		beego.Error("ParseToken:", err)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// not a token
				return nil, ErrInput
			} else if ve.Errors&
				(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// expired or active yet
				return nil, ErrExpired
			} else {
				return nil, ErrInput
			}
		} else {
			return nil, ErrInput
		}
	}

	if !token.Valid {
		beego.Error("TokenInvalid:", tokenString)
		return nil, ErrInput
	}
	beego.Debug("Token:", token)
	return token, nil
}
func (c *BaseController) AuthFail() {
	http.Error(c.Ctx.ResponseWriter, "Not logged in", http.StatusUnauthorized)
}
