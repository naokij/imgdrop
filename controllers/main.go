package controllers

type MainController struct {
	BaseController
}

func (this *MainController) Get() {
	this.CheckLogin()
	this.SetPagetitle(this.Tr("nav_upload"))
	this.Layout = "layout.html"
	this.TplName = "index.html"
}
