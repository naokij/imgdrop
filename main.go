/*
Copyright 2016 Jiang Le

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

package main

import (
	"github.com/astaxie/beego"
	"github.com/naokij/imgdrop/controllers"
	"github.com/naokij/imgdrop/routes"
	"github.com/naokij/imgdrop/setting"
)

func init() {

}

func main() {
	setting.InitApp()
	routes.InitRoutes()
	controllers.InitLocales()
	beego.SetStaticPath("/static", "static")
	beego.SetStaticPath("/images", "images")
	beego.Run()
}
