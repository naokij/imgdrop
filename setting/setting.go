/*
Copyright 2014 Jiang Le

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

package setting

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	AppName       string
	AppHost       string
	AppUrl        string
	AppLogo       string
	TmpPath       string
	ImageProxyURL string
)

var (
	CookieUserName     string
	CookieRememberName string
)

var (
	ConfigBroken bool
)

var (
	Cache cache.Cache
)

const (
	AppVer = "VERSION 0.0.1"
)

func InitApp() {
	var err error

	AppName = beego.AppConfig.String("appname")
	AppHost = beego.AppConfig.String("apphost")
	AppUrl = beego.AppConfig.String("appurl")
	AppLogo = beego.AppConfig.String("applogo")
	TmpPath = beego.AppConfig.String("tmppath")
	CookieUserName = beego.AppConfig.String("cookieusername")
	CookieRememberName = beego.AppConfig.String("CookieRememberName")
	imageproxyurl, _ := beego.GetConfig("String", "imageproxyurl", "")
	ImageProxyURL = imageproxyurl.(string)

	if err = orm.RegisterDataBase("default", "sqlite3", "data.db", 30); err != nil {
		beego.Error("sqlite3", err)
		ConfigBroken = true
	}

	// cache system
	Cache, err = cache.NewCache("file", `{"CachePath":"./cache","FileSuffix":".cache","DirectoryLevel":2,"EmbedExpiry":120}`)
	if err != nil {
		beego.Error("cache", err)
		ConfigBroken = true
	}

}
