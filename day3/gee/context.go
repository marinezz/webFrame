package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 重名了H，方便构造json数据
type H map[string]interface{}

// Context context 第一步：构造属性
type Context struct {
	// 请求与响应
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求路径、方法、参数
	Path   string
	Method string
	Params map[string]string
	// 响应码
	StatusCode int
}

// 第二部：构造函数
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// Query 第三步 ：获取参数的Query，PostForm，Param方法
// /hello?name=marine
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// PostForm post 请求体中的参数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Param param 比如 /hello/name/marine  name = marine 就是参数
func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) setHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) setCode(code int) {
	// 将请求放入响应中
	c.Writer.WriteHeader(code)
	// 放入context中
	c.StatusCode = code
}

// 第四步：实现String、JSON、Data、HTML快速构建方法
// String 用于构建错误信息
func (c *Context) String(code int, format string, value ...interface{}) {
	// 构造响应头
	c.setHeader("Content-Type", "text-plain")
	// 构造状态码
	c.setCode(code)
	// 构造响应体
	// fmt.Sprintf 格式化并返回字符串
	c.Writer.Write([]byte(fmt.Sprintf(format, value)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.setHeader("Content-Tpye", "application/json")
	c.setCode(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.setCode(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.setHeader("Content-Type", "text/html")
	c.setCode(code)
	c.Writer.Write([]byte(html))
}
