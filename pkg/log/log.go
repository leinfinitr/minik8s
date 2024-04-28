package log

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
)

// 日志级别对应的颜色分别是
// Info: 绿色
// Error: 红色
// Warn: 黄色
// Debug: 蓝色

var ifDebug = true
var ifPrint = true
var logPath = "./k8s.log"
var File *os.File = nil

func init() {
	// 创建一个 k8s.log 文件
	var err error
	File, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("打开日志文件失败")
		return
	}
}

// WriteLogToFile 将日志写入到文件中
func WriteLogToFile(component string, msg string) string {
	// 1. 获取当前的时间
	t := time.Now()
	currentTimeStr := t.Format("2006-01-02 15:04:05")
	// 2. 获取发起调用的函数名，文件名，行号
	funcName, file, line, _ := runtime.Caller(2)
	// 3. 组装成字符串
	logStr := fmt.Sprintf("[%s] %s %s %s %d: %s\n", currentTimeStr, file, component, runtime.FuncForPC(funcName).Name(), line, msg)
	// 4. 将字符串写入到文件中
	File.WriteString(logStr)
	return logStr
}

// InfoLog 基本的日志信息
func InfoLog(component string, msg string) {
	logStr := WriteLogToFile(component, msg)
	if ifPrint {
		color.Green(logStr)
	}
}

// DebugLog Debug日志
func DebugLog(component string, msg string) {
	if ifDebug {
		logStr := WriteLogToFile(component, msg)
		if ifPrint {
			color.Blue(logStr)
		}
	}
}

// WarnLog 警告日志
func WarnLog(component string, msg string) {
	logStr := WriteLogToFile(component, msg)
	if ifPrint {
		color.Yellow(logStr)
	}
}

// ErrorLog 错误日志
func ErrorLog(component string, msg string) {
	logStr := WriteLogToFile(component, msg)
	if ifPrint {
		color.Red(logStr)
	}
}
