package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

func V2Log() gee.HandleFunc {
	return func(c *gee.Context) {
		now := time.Now()

		// c.Fail(500, "Internal Server Error")

		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(now))
	}
}

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	r.GET("/hello", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1> hello golang</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(V2Log())
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s,you,re at %s\n", c.Param("name"), c.Path)
		})
	}
	r.Run(":9090")
}
