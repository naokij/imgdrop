package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/validation"
	"github.com/naokij/imgdrop/models"
	"github.com/naokij/imgdrop/setting"
	"net/url"
)

//登录表单
type LoginForm struct {
	Username string `form:"Username,text"valid:"Required;"`
	Password string `form:"Password,password"valid:"Required;"`
	Remember string `form:"Remember,text"`
}

type AuthController struct {
	BaseController
}

func (this *AuthController) Login() {
	if this.IsLogin {
		this.Redirect("/", 302)
	}
	returnUrl := this.Ctx.Input.Param(":returnurl")
	if returnUrl != "" {
		u, err := url.Parse(returnUrl)
		if err == nil {
			if u.Host == setting.AppHost {
				this.SetSession("ReturnUrl", returnUrl)
			}
		}
	}
	form := LoginForm{}
	valid := validation.Validation{}
	this.loginPage(&form, &valid)
}

//显示登录页面
func (this *AuthController) loginPage(form *LoginForm, valid *validation.Validation) {
	this.SetPagetitle(this.Tr("login"))
	//this.Layout = "layout.html"
	this.TplName = "login.html"
	this.Data["form"] = form
	this.Data["errors"] = valid.Errors
	this.Data["errorsMap"] = valid.ErrorMap()
	this.Data["HasError"] = valid.HasErrors()
}
func (this *AuthController) DoLogin() {
	if this.IsLogin {
		this.Redirect("/", 302)
	}
	this.CheckOnceToken()
	valid := validation.Validation{}
	form := LoginForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Error(err)
	}
	b, err := valid.Valid(form)
	if err != nil {
		beego.Error(err)
	}
	if !b {
		this.loginPage(&form, &valid)
		return
	}
	//用户不存在？
	user := models.User{Username: form.Username, Email: form.Username}
	if err := user.Read("Username"); err != nil {
		if err2 := user.Read("Email"); err2 != nil {
			errMsg := fmt.Sprintf("用户 %s 不存在!", form.Username)
			beego.Trace(errMsg)
			valid.SetError("Username", errMsg)
			this.loginPage(&form, &valid)
			return
		}
	}

	//检查密码
	if !user.VerifyPassword(form.Password) {
		beego.Trace(fmt.Sprintf("%s 登录失败！", form.Username))
		valid.SetError("Password", "密码错误")
		this.loginPage(&form, &valid)
		return
	}
	//验证全部通过
	var remember bool
	if form.Remember != "" {
		remember = true
	}
	this.LogUserIn(&user, remember)
	this.Redirect(GetLoginRedirectUrl(this.Ctx), 302)
	return
}

func GetLoginRedirectUrl(ctx *context.Context) (returnUrl string) {
	var ok bool
	if returnUrl, ok = ctx.Input.CruSession.Get("ReturnUrl").(string); returnUrl != "" && ok {
		ctx.Input.CruSession.Delete("ReturnUrl")
	} else {
		returnUrl = "/"
	}
	return returnUrl
}

func (this *AuthController) Logout() {
	this.LogUserOut()
	this.Redirect("/", 302)
	return
}
