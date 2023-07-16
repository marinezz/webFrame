package main

import (
	"fmt"
	"net/http"
)

// 手动实现handler，只要传入了一个ServerHTTP接口实例，所有的HTTP请求都交给该实例
//type Handler interface {
//	ServeHTTP(ResponseWriter, *Request)
//}

// Engine 建一个空结构体
type Engine struct{}

// 实现ServeHTTP
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 拿到请求路径
	path := req.URL.Path
	switch path {
	case "/":
		fmt.Fprintf(w, "URL.path:%v\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Hander[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404")
	}
}

// 和base1中的实现是一样的，按照不同的实现方法
func main() {
	// 两个参数第一个是端口，第二个是一个handler接口
	// go语言中会将结构体强制转换成接口，所以传递过去会调用handler.ServerHTTP()方法，上面已经实现
	http.ListenAndServe(":9090", new(Engine))
}
