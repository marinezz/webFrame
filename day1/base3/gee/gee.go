package gee

import (
	"fmt"
	"net/http"
)

// 本质是构造一个路由表，serveHTTP根据路由表，找到对应的处理函数，将请求和响应都传递给处理函数
// 分为以下几步
// 1、 定义结构体 ，一个路由map:key为[方法 + 路径]，value为具体的执行函数
// 2、 实现构造函数
// 3、 依次实现GET和POST方法
// 4、 实现ServeHTTP，根据路由表找到执行函数
// 5、 实现Run()方法，启动服务

// HandlerFunc 类型别名
type HandlerFunc func(w http.ResponseWriter, req *http.Request)

// Engine 第一步：定义结构体
type Engine struct {
	router map[string]HandlerFunc
}

// New 第二步：实现New()方法，作为构造函数
func New() *Engine {
	// make 内存分配
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 拼接key，并加入路由表（首字母小写，不对外暴露）
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	// 拼接key 方法 + 路径
	key := method + pattern
	// 加入路由表
	engine.router[key] = handler
}

// GET 第三步：实现GET和POST方法(传入的参数是路径和执行方法)(就是将路径和方法放入路由表中)
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRouter("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, hanlder HandlerFunc) {
	engine.addRouter("POST", pattern, hanlder)
}

// ServeHTTP 第四步：实现ServerHTTP方法（根据路由表解析出方法）
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 根据请求的方法和路径拼接出Key
	key := req.Method + req.URL.Path
	// 找到方法，将请求和响应传递给方法
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		// 没有找到方法，就直接报错
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

// Run 第五步，实现run方法，启动服务，本质上就是调用
func (engine *Engine) Run(addr string) error {
	// 传入端口号和
	return http.ListenAndServe(addr, engine)
}
