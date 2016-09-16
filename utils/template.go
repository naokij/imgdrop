package utils

import (
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"github.com/naokij/imgdrop/setting"
	"github.com/rs/xid"
	"html/template"
	"reflect"
	"strings"
	"time"
)

func loadtimes(t time.Time) int {
	return int(time.Since(t).Nanoseconds() / 1e6)
}

func nl2br(text string) template.HTML {
	return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
}

//生成防止重复提交token
func OnceToken() (token string) {
	guid := xid.New()
	token = guid.String()
	if err := setting.Cache.Put("Once_"+token, 1, 86400); err != nil {
		beego.Error("cache", err)
	}
	beego.Trace("Once-test", token)
	return token
}

func OnceFormHtml() template.HTML {
	return template.HTML("<input type=\"hidden\" name=\"_once\" value=\"" +
		OnceToken() + "\"/>")
}

func LoginUrlFor(endpoint string, values ...interface{}) string {
	return beego.URLFor("AuthController.Login", ":returnurl", template.URLQueryEscaper(beego.URLFor(endpoint, values...)))
}

// isSet returns whether a given array, channel, slice, or map has a key
// defined.
func isSet(a interface{}, key interface{}) bool {
	av := reflect.ValueOf(a)
	kv := reflect.ValueOf(key)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		if int64(av.Len()) > kv.Int() {
			return true
		}
	case reflect.Map:
		if kv.Type() == av.Type().Key() {
			return av.MapIndex(kv).IsValid()
		}
	}

	return false
}

func init() {
	// Register template functions.
	beego.AddFuncMap("loadtimes", loadtimes)
	beego.AddFuncMap("jsescape", template.JSEscapeString)
	beego.AddFuncMap("nl2br", nl2br)
	beego.AddFuncMap("oncetoken", OnceToken)
	beego.AddFuncMap("onceformhtml", OnceFormHtml)
	beego.AddFuncMap("loginurl", LoginUrlFor)
	beego.AddFuncMap("i18n", i18n.Tr)
	beego.AddFuncMap("isset", isSet)
}
