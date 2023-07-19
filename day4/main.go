package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/hello", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Group("/marine")
	{
		v1.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "name: marine,path %s", c.Path)
		})

		v1.GET("/hello/:name", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"name": c.Param("name"),
				"path": c.Path,
			})
		})
	}

	r.Run(":9090")
}
