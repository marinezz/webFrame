package main

import (
	"fmt"
	"net/http"
)

func main() {
	//  原生写法，返回一个路径，写入具体的路径
	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(writer, "URL.path:%v\n", req.URL.Path)
	})

	// 另外一种写法，处理函数不用隐函数，写在外面
	http.HandleFunc("/hello", hellohander)
	http.ListenAndServe(":9090", nil)
}

func hellohander(writer http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(writer, "Hander[%q] = %q\n", k, v)
	}
}
