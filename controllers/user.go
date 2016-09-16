package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/naokij/imgdrop/models"
	"strconv"
)

type ProfileForm struct {
	Password        string
	PasswordConfirm string
	Email           string `valid:"Email"`
}

type UserForm struct {
	Username string
	Email    string `valid:"Email"`
	Password string
}

type UserController struct {
	BaseController
}

func (this *UserController) Get() {
	this.CheckLogin()
	form := ProfileForm{}
	form.Email = this.User.Email
	form.Password = ""
	valid := validation.Validation{}
	this.profileForm(&form, &valid)
}

func (this *UserController) profileForm(form *ProfileForm, valid *validation.Validation) {
	this.SetPagetitle(this.Tr("nav_profile"))
	this.Data["form"] = form
	this.Layout = "layout.html"
	this.TplName = "profile.html"
	this.Data["errors"] = valid.Errors
	this.Data["errorsMap"] = valid.ErrorMap()
	this.Data["HasError"] = valid.HasErrors()
}

func (this *UserController) Post() {
	this.CheckLogin()
	this.CheckOnceToken()
	this.FlashRead("notice")
	valid := validation.Validation{}
	form := ProfileForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Error(err)
	}
	b, err := valid.Valid(form)
	if err != nil {
		beego.Error(err)
	}
	if !b {
		this.profileForm(&form, &valid)
		return
	}
	this.User.Email = form.Email
	if len(form.Password) > 0 {
		if form.Password != form.PasswordConfirm {
			valid.SetError("Password", this.Tr("password_confirm_error"))
			this.profileForm(&form, &valid)
			return
		}
		this.User.SetPassword(form.Password)
	}
	err = this.User.Update("Email", "Password", "Salt")
	if err != nil {
		beego.Error(err)
	}
	this.FlashWrite("notice", this.Tr("updated_message"))
	this.profileForm(&form, &valid)

}

func (this *UserController) newUserForm(form *UserForm, valid *validation.Validation) {
	this.FlashRead("success")
	this.Data["form"] = form
	this.TplName = "new_user.html"
	this.Layout = "layout.html"
	this.Data["errors"] = valid.Errors
	this.Data["errorsMap"] = valid.ErrorMap()
	this.Data["HasError"] = valid.HasErrors()

}

func (this *UserController) editUserForm(form *UserForm, valid *validation.Validation) {
	this.FlashRead("success")
	this.Data["form"] = form
	this.TplName = "edit_user.html"
	this.Layout = "layout.html"
	this.Data["errors"] = valid.Errors
	this.Data["errorsMap"] = valid.ErrorMap()
	this.Data["HasError"] = valid.HasErrors()

}

func (this *UserController) New() {
	this.CheckAdmin()
	this.SetPagetitle(this.Tr("nav_users_new"))
	this.TplName = "new_user.html"
	valid := validation.Validation{}
	form := UserForm{}
	this.newUserForm(&form, &valid)
}

func (this *UserController) Create() {
	this.CheckAdmin()
	user := models.User{}
	this.CheckOnceToken()
	this.SetPagetitle(this.Tr("nav_users_new"))
	this.TplName = "edit_user.html"
	valid := validation.Validation{}
	form := UserForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Error(err)
	}
	_, err := valid.Valid(form)
	if err != nil {
		beego.Error(err)
	}
	if len(form.Password) == 0 {
		valid.SetError("Password", this.Tr("ErrorRequired"))
	}
	if len(form.Username) == 0 {
		valid.SetError("Username", this.Tr("ErrorRequired"))
	}
	user.Username = form.Username
	if err = user.ValidUsername(); err != nil {
		valid.SetError("Username", err.Error())
	}
	testUser := models.User{Username: form.Username}
	err = testUser.Read("Username")
	if err == nil {
		if testUser.Id != user.Id {
			valid.SetError("Username", this.Tr("ErrorDuplicatedUsername"))
		}
	}
	testUser = models.User{Email: form.Email}
	err = testUser.Read("Email")
	if err == nil {
		if testUser.Id != user.Id {
			valid.SetError("Email", this.Tr("ErrorDuplicatedEmail"))
		}
	}
	if !valid.HasErrors() {
		user.Email = form.Email
		user.SetPassword(form.Password)
	} else {
		this.newUserForm(&form, &valid)
		return
	}
	err = user.Insert()
	if err != nil {
		beego.Error(err)
		return
	}
	this.FlashWrite("success", this.Tr("created_message"))
	this.newUserForm(&form, &valid)
	return
}

func (this *UserController) Edit() {
	this.CheckAdmin()
	user := this.getUserFromRequest()
	this.SetPagetitle(this.Tr("nav_users_edit"))
	this.TplName = "edit_user.html"
	valid := validation.Validation{}
	form := UserForm{
		Username: user.Username,
		Email:    user.Email,
	}
	this.editUserForm(&form, &valid)
}

func (this *UserController) Update() {
	this.CheckAdmin()
	user := this.getUserFromRequest()
	this.CheckOnceToken()
	this.SetPagetitle(this.Tr("nav_users_edit"))
	this.TplName = "edit_user.html"
	valid := validation.Validation{}
	form := UserForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Error(err)
	}
	form.Username = user.Username
	_, err := valid.Valid(form)
	if err != nil {
		beego.Error(err)
	}
	testUser := models.User{Username: form.Username}
	err = testUser.Read("Username")
	if err == nil {
		if testUser.Id != user.Id {
			valid.SetError("Username", this.Tr("ErrorDuplicatedUsername"))
		}
	}
	testUser = models.User{Email: form.Email}
	err = testUser.Read("Email")
	if err == nil {
		if testUser.Id != user.Id {
			valid.SetError("Email", this.Tr("ErrorDuplicatedEmail"))
		}
	}
	if valid.HasErrors() {
		this.editUserForm(&form, &valid)
		return
	}
	user.Email = form.Email
	if len(form.Password) > 0 {
		user.SetPassword(form.Password)
	}
	err = user.Update("Email", "Password", "Salt")
	if err != nil {
		beego.Error(err)
		return
	}
	this.FlashWrite("success", this.Tr("updated_message"))
	this.editUserForm(&form, &valid)
}

func (this *UserController) Del() {
	this.CheckAdmin()
}

func (this *UserController) List() {
	this.CheckAdmin()
	this.SetPagetitle(this.Tr("nav_manage_users"))
	var users []*models.User
	_, err := models.Users().All(&users)
	if err != nil {
		beego.Error(err)
	}
	this.Data["Users"] = users
	this.TplName = "users.html"
	this.Layout = "layout.html"
}
func (this *UserController) getUserFromRequest() models.User {
	this.Data["Id"] = this.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
	user := models.User{Id: id}
	err := user.Read()
	if err != nil {
		this.Abort("404")
	}
	return user
}
