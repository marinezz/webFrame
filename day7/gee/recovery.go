package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// 触发panic的堆栈信息
func trace(message string) string {
	// uintptr 无符号整形，可以保存一个指针地址
	var pcs [32]uintptr
	// 获取函数调用信息，跳过 3 层 ，并用pcs接收 (调用栈的程序计数器)
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(message + "\nTracebacke:")

	for _, pc := range pcs[:n] {
		// 获取程序计数器对应的函数的信息
		fn := runtime.FuncForPC(pc)
		// 返回文件名和行号
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}

	return str.String()
}

// Recovery 如果发生错误，则捕获异常
func Recovery() HandleFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				// 打印错误信息
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internet server error")
			}
		}()

		c.next()
	}
}
