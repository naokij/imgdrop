package controllers

import (
	"github.com/naokij/imgdrop/setting"
	"io/ioutil"
	"strings"
)

type ImagesController struct {
	BaseController
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
	images := []string{}
	images = make([]string, 0)

	for _, file := range files {
		fileName := file.Name()
		if strings.HasPrefix(fileName, ".") {
			continue
		}
		images = append(images, urlPrefix+file.Name())
	}

	pers := 10
	total := int64(len(images))
	pager := this.SetPaginator(pers, total)

	paged_images := []string{}
	paged_images = make([]string, 0)
	for k, url := range images {
		if k >= pager.Offset() && k < pager.Offset()+pers {
			paged_images = append(paged_images, url)
		}
	}
	this.Data["images"] = paged_images

	//this.Data["ExtraCSS"] = "/static/css/images.css"
	this.Layout = "layout.html"
	this.TplName = "images.html"
}
