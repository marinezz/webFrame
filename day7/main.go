package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.Default()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "hello golang!!")
	})

	r.GET("/panic", func(c *gee.Context) {
		names := []string{"marine"}
		// 越界
		c.String(http.StatusOK, names[10])
	})

	r.Run(":9090")
}
