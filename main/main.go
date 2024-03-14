package main

import (
	"html/template"
	"net/http"
	"router"
	"util"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"UnixToTime": util.TemplateFunc,
	})

	r.Static("/static", "static")
	r.StaticFile("/favicon.ico", "static/resources/favicon.ico")
	r.LoadHTMLGlob("templates/*")
	r.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/")
	})

	router.UserRouterInit(r)
	router.IdexRouterInit(r)

	r.Run("localhost:8080")
	// r.RunTLS("localhost:8080", "./cert/test.pem", "./cert/test.key")

}
