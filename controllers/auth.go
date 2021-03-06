package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"net/url"
	"ssoauth/models"
)

type AuthController struct {
	BaseController
}

func (c *AuthController) Get() {

	token, err := c.ParseToken("header")
	if err != nil {
		c.AuthFail()
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := claims["sub"].(string)
	beego.Debug("ParseUserId:", userId)

	user := models.User{}
	if code, err := user.FindById(userId); err != nil {
		beego.Error("FindUserById:", err)
		if code == models.ERROR_NOT_FOUND {
			c.AuthFail()
		} else {
			c.AuthFail()
		}
		return
	}

	if !user.HasDomain(c.Ctx.Input.Domain()) {
		beego.Debug("CheckApp:", c.Ctx.Input.Domain(), "Failed")
		c.AccessForbidden()
		return
	}

	userAgent := c.ParseUserAgent()

	ok := user.HasToken(userAgent, token.Raw)
	if !ok {
		beego.Debug("HasToken:", "Token is valid but revoked by user.")
		c.AuthFail()
		return
	}

	w := c.Ctx.ResponseWriter
	w.Header().Set("Igenetech-User-Id", user.Id.Hex())
	w.Header().Set("Igenetech-User-Name", user.Name)
	w.Header().Set("Igenetech-User-Email", user.Email)
	userRoleBytes, _ := json.Marshal(user.Roles)
	w.Header().Set("igenetech-User-Role", string(userRoleBytes))
	roles := ""
	for k, _ := range user.Roles {
		if k != "" {
			roles += k
			roles += ","
		}
	}

	w.Header().Set("Igenetech-User-Roles", roles)
	c.Ctx.Output.SetStatus(http.StatusOK)
	beego.Debug("AuthSuccess:", c.Ctx.Input.IP)
	return
}

type AuthByCookieController struct {
	BaseController
}

func (c *AuthByCookieController) Get() {
	token, err := c.ParseToken("cookie")
	if err != nil {
		c.AuthFail()
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["sub"].(string)
	beego.Debug("ParseUserId:", userId)

	user := models.User{}
	if code, err := user.FindById(userId); err != nil {
		beego.Error("FindUserById:", err)
		if code == models.ERROR_NOT_FOUND {
			c.AuthFail()
		} else {
			c.AuthFail()
		}
		return
	}

	userAgent := c.ParseUserAgent()

	ok := user.HasToken(userAgent, token.Raw)
	if !ok {
		beego.Debug("HasToken:", "Token is valid but revoked by user.")
		c.AuthFail()
		return
	}

	if !user.HasDomain(c.Ctx.Input.Domain()) {
		beego.Debug("CheckApp:", c.Ctx.Input.Domain(), "Failed")
		c.AccessForbidden()
		return
	}

	w := c.Ctx.ResponseWriter
	userName := ""
	encodedUserName, encodeErr := url.Parse(user.Name)
	if encodeErr != nil {
		beego.Error("UserNameEncode", encodeErr)
		userName = user.Name
	} else {
		userName = encodedUserName.EscapedPath()
	}
	w.Header().Set("Igenetech-User-Id", user.Id.Hex())
	w.Header().Set("Igenetech-User-Name", userName)
	w.Header().Set("Igenetech-User-Email", user.Email)
	userRoleBytes, _ := json.Marshal(user.Roles)
	w.Header().Set("igenetech-User-Role", string(userRoleBytes))
	roles := ""
	for k, _ := range user.Roles {
		if k != "" {
			roles += k
			roles += ","
		}
	}

	w.Header().Set("Igenetech-User-Roles", roles)
	beego.Debug("Header:", w.Header())
	c.Ctx.Output.SetStatus(http.StatusOK)
	beego.Debug("AuthSuccess:", c.Ctx.Input.IP)
	return
}
