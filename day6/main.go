package main

import (
	"fmt"
	"gee"
	"html/template"
	"net/http"
	"time"
)

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("template/*")
	r.Static("/assets", "./static")

	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.Run(":9090")
}
