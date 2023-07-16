package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 1、 构建context结构体，包含属性有：请求、响应、请求路径、请求方法、响应状态码
// 2、 构造函数
// 3、 Query方法和PostForm方法，查询参数信息
// 4、 提供快速构造String，Data，JSON，HTML的方法() ,分三部，构造响应头，构造状态码，构造响应体

// H 为map[string]interface{}取了个别名为H,构造JSON数据的时候更加便捷
type H map[string]interface{}

// Context 第一步：构造结构体
type Context struct {
	// 请求与响应
	Writer http.ResponseWriter
	Req    *http.Request

	// 请求路径与方法
	Path   string
	Method string

	// 响应状态码
	StatusCode int
}

// 第二步 构造函数
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// Query 第三步 实现Query方法和PostForm方法
func (c *Context) Query(key string) string {
	// 从请求路径中根据key获取参数
	return c.Req.URL.Query().Get(key)
}

func (c *Context) PostForm(key string) string {
	// 根据key获取表单数据
	return c.Req.FormValue(key)
}

// 构造响应头
func (c *Context) setHeader(key string, value string) {
	// 响应头中放入响应名称
	c.Writer.Header().Set(key, value)
}

// 构造状态码,将状态码保存至context并写入响应头
func (c *Context) setStatus(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 第四步 实现String、Data、JSON、HTML快速构建方法
// ...interface{}变长参数，会被转换成切片
// 正确调用顺序是Header.Set() 然后 WriteHeader() 然后Write()
func (c *Context) String(code int, format string, values ...interface{}) {
	// 构造响应头
	c.setHeader("Content-Type", "text/plain")
	// 构造状态码
	c.setStatus(code)
	// 构造响应体，传入一个切片
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.setHeader("Content-Type", "application/json")
	c.setStatus(code)
	// 根据io.writer创建编码器encoder，然后编码器调用encode将对象编码成JSON
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		// 编码失败则报错
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.setStatus(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.setHeader("Content-Type", "text/html")
	c.setStatus(code)
	c.Writer.Write([]byte(html))
}
