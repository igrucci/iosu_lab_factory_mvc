package render

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"html/template"
)

func RenderHtml(ctx *fasthttp.RequestCtx, fileName string, data interface{}) {
	tmpl, err := template.ParseFiles(fileName)
	if err != nil {
		logrus.Error("Error parsing template: ", err)
		return
	}
	ctx.Response.Header.SetContentType("text/html; charset=utf-8")
	err = tmpl.Execute(ctx, data)
	if err != nil {
		fmt.Println(err)
	}
}
