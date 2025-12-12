// Package tools 日志工具
package tools

import "fmt"

// Logf 暂定的一个日志方法，这里先使用一个打印语句
func Logf(format string, v ...interface{}) {
	fmt.Printf("[LOG] "+format, v...)
}
