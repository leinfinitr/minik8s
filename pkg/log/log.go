// 描述：对 println 的封装，用于控制是否输出日志

package log

import "fmt"

var log = true

func SetLog(l bool) {
	log = l
}

func Println(v ...interface{}) {
	if log {
		println(fmt.Sprint(v...))
	}
}
