package controllers

import (
	"github.com/naokij/imgdrop/setting"
	"github.com/rs/xid"
	"os"
)

type UploadController struct {
	BaseController
}

type UploadResponse struct {
	Success bool
	URL     string
	Error   string
}

func (this *UploadController) Post() {
	response := UploadResponse{}
	this.CheckLogin()
	f, h, err := this.GetFile("files[]")
	defer f.Close()
	if err != nil {
		response.Success = false
		response.Error = err.Error()
		this.Data["json"] = response
	} else {
		savePath := "./images/" + this.User.Username
		urlPrefix := "/images/" + this.User.Username + "/"
		saveFilename := savePath + "/" + h.Filename
		err = os.MkdirAll(savePath, 0755)
		if err != nil {
			response.Success = false
			response.Error = err.Error()
		}
		if _, err := os.Stat(saveFilename); err == nil {
			//duplicated file
			guid := xid.New()
			saveFilename = savePath + "/" + guid.String() + h.Filename
		}
		err = this.SaveToFile("files[]", saveFilename)
		if err != nil {
			response.Success = false
			response.Error = err.Error()
		} else {
			response.Success = true
			response.URL = setting.AppUrl + urlPrefix + h.Filename
		}

	}
	this.Data["json"] = response
	this.ServeJSON()
}
