package klog

import (
	"os"

	"github.com/fatih/color"
	// "log"
)

var logPath = ""
var logFile *os.File = nil

var outputInfoLog = true

// 最简单版本的写入log文件并输出的log，根据不同的情况进行输出
func initLog() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	logFile, err = os.Create(workDir + "/k8s.log")
	if err != nil {
		panic(err)
	}

	logPath = workDir + "/k8s.log"
}

func WarnLog(component string, logInfo string) {
	if logFile == nil {
		initLog()
	}
	// log.
	var result = "[" + component + "] : " + logInfo + "\n"
	logFile.WriteString(result)
	if outputInfoLog {
		color.Red(result)
	}

}

func InfoLog(component string, logInfo string) {
	if logFile == nil {
		initLog()
	}
	// log.
	var result = "[" + component + "] : " + logInfo + "\n"
	logFile.WriteString(result)
	if outputInfoLog {
		color.Green(result)
	}

}
