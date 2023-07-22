package gee

import (
	"log"
	"time"
)

func Logger() HandleFunc {
	return func(c *Context) {
		now := time.Now()

		c.next()

		// 打印 状态码 路径 请求响应时间
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(now))
	}
}
