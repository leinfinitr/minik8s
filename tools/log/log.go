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

var ifDebug = false
var ifFile = false
var ifPrint = true
var logPath = "./k8s.log"
var File *os.File = nil

func init() {
	// 创建一个 k8s.log 文件
	var err error
	File, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Tools log init fail!")
		return
	}
}

// WriteLogToFile 将日志写入到文件中
func WriteLogToFile(msg string) string {
	// 1. 获取当前的时间
	t := time.Now()
	currentTimeStr := t.Format("2006-01-02 15:04:05")
	// 2. 获取发起调用的函数名，文件名，行号
	funcName, file, line, _ := runtime.Caller(2)
	// 3. 组装成字符串
	logStr := fmt.Sprintf("[%s] [%s %d] [%s]: %s\n", currentTimeStr, file, line, runtime.FuncForPC(funcName).Name(), msg)
	// 4. 将字符串写入到文件中
	if ifFile {
		_, err := File.WriteString(logStr)
		if err != nil {
			return ""
		}
	}
	return logStr
}

// InfoLog 基本的日志信息
func InfoLog(msg string) {
	logStr := WriteLogToFile(msg)
	if ifPrint {
		color.Green(logStr)
	}
}

// DebugLog Debug日志
func DebugLog(msg string) {
	if ifDebug {
		logStr := WriteLogToFile(msg)
		if ifPrint {
			color.Blue(logStr)
		}
	}
}

// WarnLog 警告日志
func WarnLog(msg string) {
	logStr := WriteLogToFile(msg)
	if ifPrint {
		color.Yellow(logStr)
	}
}

// ErrorLog 错误日志
func ErrorLog(msg string) {
	logStr := WriteLogToFile(msg)
	if ifPrint {
		color.Red(logStr)
	}
}
