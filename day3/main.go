package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()

	r.GET("/hello", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>hello golang</h1>")
	})

	r.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s,path : %s\n", c.Param("name"), c.Path)
	})

	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s,path :%s\n", c.Query("name"), c.Path)
	})

	r.GET("assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"filepath": c.Param("filepath"),
			"path":     c.Path,
		})
	})

	r.GET("/hello/:name/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"path": c.Path,
			"name": c.Param("name"),
			"file": c.Param("filepath"),
		})
	})
	r.Run(":9090")
}
