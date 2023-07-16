package main

import (
	"fmt"
	"gee"
	"net/http"
)

func main() {
	// build constraints exclude all the go files in "src文件"
	// 遇到上面的问题，导致自定义包无法导入，解决方法是自己在src上创建了一个名称叫gee的空包（为什么原因未知）
	// 找到这个空包，删除即可，也要确保mod 文件的正确
	// 创建实例
	r := gee.New()

	// GET
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.path:%v\n", req.URL.Path)
	})

	// 和前面两个例子一样的
	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Hander[%q] = %q\n", k, v)
		}
	})

	// 启动服务
	r.Run(":9090")
}
