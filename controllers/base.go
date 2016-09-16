/*
Copyright 2014, 2016 Jiang Le

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/beego/i18n"
	"github.com/naokij/imgdrop/models"
	"github.com/naokij/imgdrop/setting"
	"github.com/naokij/imgdrop/utils"
	"strings"
	"time"
)

var langTypes []*langType // Languages are supported.

// langType represents a language type.
type langType struct {
	Lang, Name string
}

type BaseController struct {
	beego.Controller
	User    *models.User
	IsLogin bool
	i18n.Locale
}

var validationMessages map[string]string

func InitLocales() {
	// Initialized language type list.
	langs := strings.Split(beego.AppConfig.String("lang::types"), "|")
	names := strings.Split(beego.AppConfig.String("lang::names"), "|")
	langTypes = make([]*langType, 0, len(langs))
	for i, v := range langs {
		langTypes = append(langTypes, &langType{
			Lang: v,
			Name: names[i],
		})
	}

	for _, lang := range langs {
		beego.Trace("Loading language: " + lang)
		if err := i18n.SetMessage(lang, "conf/"+"locale_"+lang+".ini"); err != nil {
			beego.Error("Fail to set message file: " + err.Error())
			return
		}
	}
}

// setLangVer sets site language version.
func (this *BaseController) setLangVer() bool {
	isNeedRedir := false
	hasCookie := false

	// 1. Check URL arguments.
	lang := this.Input().Get("lang")

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		lang = this.Ctx.GetCookie("lang")
		hasCookie = true
	} else {
		isNeedRedir = true
	}

	// Check again in case someone modify by purpose.
	if !i18n.IsExist(lang) {
		lang = ""
		isNeedRedir = false
		hasCookie = false
	}

	// 3. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := this.Ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // Only compare first 5 letters.
			if i18n.IsExist(al) {
				lang = al
			}
		}
	}

	// 4. Default language is English.
	if len(lang) == 0 {
		lang = "en-US"
		isNeedRedir = false
	}

	curLang := langType{
		Lang: lang,
	}

	// Save language information in cookies.
	if !hasCookie {
		this.Ctx.SetCookie("lang", curLang.Lang, 1<<31-1, "/")
	}

	restLangs := make([]*langType, 0, len(langTypes)-1)
	for _, v := range langTypes {
		if lang != v.Lang {
			restLangs = append(restLangs, v)
		} else {
			curLang.Name = v.Name
		}
	}

	// Set language properties.
	this.Lang = lang
	this.Data["Lang"] = curLang.Lang
	this.Data["CurLang"] = curLang.Name
	this.Data["RestLangs"] = restLangs

	// set validation messages
	validationMessages = map[string]string{
		"Required":     this.Tr("ErrorRequired"),
		"Min":          this.Tr("ErrorMin", "%d"),
		"Max":          this.Tr("ErrorMax", "%d"),
		"Range":        this.Tr("ErrorRange", "%d", "%d"),
		"MinSize":      this.Tr("ErrorMinSize", "%d"),
		"MaxSize":      this.Tr("ErrorMaxSize", "%d"),
		"Length":       this.Tr("ErrorLength", "%d"),
		"Alpha":        this.Tr("ErrorAlpha"),
		"Numeric":      this.Tr("ErrorNumeric"),
		"AlphaNumeric": this.Tr("ErrorAlphaNumeric"),
		"Match":        this.Tr("ErrorMatch", "%s"),
		"NoMatch":      this.Tr("ErrorNoMatch", "%s"),
		"AlphaDash":    this.Tr("ErrorAlphaDash"),
		"Email":        this.Tr("ErrorEmail"),
		"IP":           this.Tr("ErrorIP"),
		"Base64":       this.Tr("ErrorBase64"),
		"Mobile":       this.Tr("ErrorMobile"),
		"Tel":          this.Tr("ErrorTel"),
		"Phone":        this.Tr("ErrorPhone"),
		"ZipCode":      this.Tr("ErrorZipCode"),
	}
	validation.SetDefaultMessage(validationMessages)
	return isNeedRedir
}

//通过session获取登录信息，并且登录
func (this *BaseController) loginViaSession() bool {
	if username, ok := this.GetSession("AuthUsername").(string); username != "" && ok {
		//beego.Trace("loginViaSession pass 1 Session[AuthUsername]" + username)
		user := models.User{Username: username}
		if user.Read("Username") == nil {
			this.User = &user
			//beego.Trace("loginViaSession pass 2 ")
			return true
		}
		beego.Trace("loginViaSession pass 2 failed ")
	}
	//beego.Trace("loginViaSession failed ")
	return false
}

//通过remember cookie获取登录信息，并且登录
func (this *BaseController) loginViaRememberCookie() (success bool) {
	username := this.Ctx.GetCookie(setting.CookieUserName)
	if len(username) == 0 {
		return false
	}

	defer func() {
		if !success {
			this.DeleteRememberCookie()
		}
	}()

	user := models.User{Username: username}
	if err := user.Read("Username"); err != nil {
		return false
	}

	secret := utils.EncodeMd5(user.Salt + user.Password)
	value, _ := this.Ctx.GetSecureCookie(secret, setting.CookieRememberName)
	if value != username {
		return false
	}
	this.User = &user
	this.LogUserIn(&user, true)

	return true
}

//删除记忆登录cookie
func (this *BaseController) DeleteRememberCookie() {
	this.Ctx.SetCookie(setting.CookieUserName, "", -1)
	this.Ctx.SetCookie(setting.CookieRememberName, "", -1)
}

//登录用户
func (this *BaseController) LogUserIn(user *models.User, remember bool) {
	this.SessionRegenerateID()
	this.SetSession("AuthUsername", user.Username)
	if remember {
		secret := utils.EncodeMd5(user.Salt + user.Password)
		days := 86400 * 30
		this.Ctx.SetCookie(setting.CookieUserName, user.Username, days)
		this.SetSecureCookie(secret, setting.CookieRememberName, user.Username, days)
	}
}

//登出用户
func (this *BaseController) LogUserOut() {
	this.DeleteRememberCookie()
	this.DelSession("AuthUsername")
	this.DestroySession()
}

func (this *BaseController) Prepare() {
	if setting.ConfigBroken {
		this.Abort("500")
	}

	// page start time
	this.Data["PageStartTime"] = time.Now()
	this.Data["AppName"] = setting.AppName
	this.Data["AppVer"] = setting.AppVer
	this.Data["PageTitle"] = setting.AppName

	// start session
	this.StartSession()

	// Redirect to make URL clean.
	if this.setLangVer() {
		i := strings.Index(this.Ctx.Request.RequestURI, "?")
		this.Redirect(this.Ctx.Request.RequestURI[:i], 302)
		return
	}

	//从session中读取登录信息
	switch {
	// save logined user if exist in session
	case this.loginViaSession():
		this.IsLogin = true
	// save logined user if exist in remember cookie
	case this.loginViaRememberCookie():
		this.IsLogin = true
	}

	if this.IsLogin {
		this.Data["User"] = &this.User
		this.Data["IsLogin"] = this.IsLogin

	}

	// read flash message
	beego.ReadFromRequest(&this.Controller)

}

// read beego flash message
func (this *BaseController) FlashRead(key string) (string, bool) {
	if data, ok := this.Data["flash"].(map[string]string); ok {
		value, ok := data[key]
		return value, ok
	}
	return "", false
}

// write beego flash message
func (this *BaseController) FlashWrite(key string, value string) {
	flash := beego.NewFlash()
	flash.Data[key] = value
	flash.Store(&this.Controller)
}

//验证防重复提交token
func (this *BaseController) CheckOnceToken() {
	token := this.GetString("_once")
	cache_key := "Once_" + token
	if token == "" {
		this.Abort("Once")
	}
	beego.Trace("cache_key", cache_key)
	beego.Trace("cache", setting.Cache.Get(cache_key))
	if setting.Cache.IsExist(cache_key) {
		setting.Cache.Delete(cache_key)
	} else {
		this.Abort("Once")
	}
}

func (this *BaseController) SetPaginator(per int, nums int64) *utils.Paginator {
	p := utils.NewPaginator(this.Ctx.Request, per, nums)
	this.Data["paginator"] = p
	return p
}

func (this *BaseController) CheckLogin() {
	if !this.IsLogin {
		this.Redirect("/login", 302)
		return
	}
}
func (this *BaseController) CheckAdmin() {
	this.CheckLogin()
	if this.User.Username != "admin" {
		this.Abort("403")
		return
	}
}

func (this *BaseController) SetPagetitle(title string) {
	this.Data["PageTitle"] = title + " | " + setting.AppName
}
