package gee

import (
	"log"
	"net/http"
)

// 将路由的结构和方法放在此文件中
// 主要的两个方法有：加入路由表和根据路由表拿到处理方法

// 1、构造路由表结构体
// 2、构造函数
// 3、实现addRoute（添加路由记录）
// 4、根据路由表找出处理函数

// 第一步：路由表结构体
type router struct {
	handlers map[string]HandleFunc
}

// 第二部：构造函数
func newRouter() *router {
	return &router{handlers: map[string]HandleFunc{}}
}

// 第三步：加入路由
func (r *router) addRouter(method string, pattern string, handler HandleFunc) {
	// 打印日志
	log.Printf("Route %4s - %4s", method, pattern)
	// 构造键
	key := method + "-" + pattern
	// 加入路由表
	r.handlers[key] = handler
}

// 第四步：拿到处理方法（放入context中）
func (r *router) handle(c *Context) {
	// 拼接key
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		// 拿到指向方法，将context放如执行方法中
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND", c.Path)
	}

}
