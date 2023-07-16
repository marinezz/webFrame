package gee

import "net/http"

// 框架的主入口
// 1、 定义结构体
// 2、 构造方法
// 3、实现GET和POST方法
// 4、重写ServeHTTP()方法
// 5、实现run方法

// HandleFunc 处理函数
type HandleFunc func(*Context)

// Engine 第一步：定义结构体
type Engine struct {
	router *router
}

// New 第二部 ：构造方法
func New() *Engine {
	return &Engine{router: newRouter()}
}

// 加入路由表
func (engine *Engine) addRouter(method string, pattern string, handler HandleFunc) {
	engine.router.addRouter(method, pattern, handler)
}

// GET 第三步：实现GET和POST方法
func (engine *Engine) GET(pattern string, handler HandleFunc) {
	// 加入路由
	engine.addRouter("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandleFunc) {
	engine.addRouter("POST", pattern, handler)
}

// 第四步：重写ServeHTTP方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}

// Run 第五步：实现Run方法
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
