package gee

import (
	"fmt"
	"reflect"
	"testing"
)

// 单元测试
func newTestRouter() *Router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello", nil)
	r.addRoute("GET", "/hello/:name", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	//  reflect.DeepEqual 判断两个值是否一致
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, params := r.getRouter("GET", "/hello/marine")

	fmt.Printf("path: %s, parmas:%s\n", n.pattern, params["name"])
}
