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

package controllers

import (
	"github.com/astaxie/beego"
	"github.com/naokij/imgdrop/setting"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error404() {
	c.Data["Title"] = "页面未找到 | " + setting.AppName
	c.TplName = "errors/404.html"
}

func (c *ErrorController) Error403() {
	c.Data["Title"] = "禁止访问 | " + setting.AppName
	c.TplName = "errors/403.html"
}

func (c *ErrorController) Error500() {
	c.Data["Title"] = "服务器内部错误 | " + setting.AppName
	c.TplName = "errors/50x.html"
}

func (c *ErrorController) Error501() {
	c.Error500()
}

func (c *ErrorController) ErrorOnce() {
	c.Data["Title"] = "重复提交 | " + setting.AppName
	c.TplName = "errors/once.html"
}
