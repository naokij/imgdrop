package controllers

import (
	"github.com/astaxie/beego"
	"github.com/naokij/imgdrop/setting"
	"io/ioutil"
	"os"
	"strings"
)

type ImagesController struct {
	BaseController
}

type Image struct {
	URL  string
	File string
}

func (this *ImagesController) Get() {
	this.CheckLogin()
	this.FlashRead("error")

	pathPrefix := "./images/" + this.User.Username + "/"
	urlPrefix := setting.AppUrl + "/images/" + this.User.Username + "/"
	files, err := ioutil.ReadDir(pathPrefix)
	if err != nil {
		this.FlashWrite("error", this.Tr("no_image_found"))
	}
	var images []Image
	images = make([]Image, 0)

	for _, file := range files {
		fileName := file.Name()
		if strings.HasPrefix(fileName, ".") {
			continue
		}
		image := Image{URL: urlPrefix + file.Name(), File: file.Name()}
		images = append(images, image)
	}

	pers := 10
	total := int64(len(images))
	pager := this.SetPaginator(pers, total)

	paged_images := []Image{}
	paged_images = make([]Image, 0)
	for k, image := range images {
		if k >= pager.Offset() && k < pager.Offset()+pers {
			paged_images = append(paged_images, image)
		}
	}
	this.Data["images"] = paged_images

	//this.Data["ExtraCSS"] = "/static/css/images.css"
	this.Layout = "layout.html"
	this.TplName = "images.html"
}

func (this *ImagesController) Delete() {
	this.CheckLogin()
	pathPrefix := "./images/" + this.User.Username + "/"
	image := this.GetString("image")
	if len(image) == 0 {
		this.Abort("500")
		return
	}
	err := os.Remove(pathPrefix + image)
	if err != nil {
		beego.Trace("Failed to remove " + pathPrefix + image)
		this.FlashWrite("error", this.Tr("image_remove_failed"))
		this.Redirect("/my_images", 302)
		return
	}
	this.FlashWrite("success", this.Tr("image_removed"))
	this.Redirect("/my_images", 302)
}
