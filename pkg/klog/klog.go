package klog

import (
	"fmt"
	"os"
	// "log"
)

var debug = false
var logPath = ""
var logFile *os.File = nil

// 最简单版本的写入log文件并输出的log，之后再进行更新
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

func appendLog(lofInfo string) {
	if logFile == nil {
		initLog()
	}
	// log.
	logFile.WriteString(lofInfo)
	fmt.Println(lofInfo)

}
