package handler

import (
	"github.com/gin-gonic/gin"
	"minik8s/tools/log"
	"strconv"
	"strings"
)

// RunServerlessWorkflow 运行Serverless Workflow
func RunServerlessWorkflow(c *gin.Context) {
	log.DebugLog("RunServerlessWorkflow")
	// 从请求中获取 serverless workflow 内容
	var content = ""
	err := c.ShouldBindJSON(&content)
	if err != nil {
		log.ErrorLog("RunServerlessWorkflow: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	param := c.Param("param")
	// 将 workflow 内容按行分割
	workflow := strings.Split(content, "\n")

	// 循环解析 workflow 内容，并根据关键字类型执行对应操作
	// 每个 function 将上一个 function 的结果作为输入，并一定会有一个输出
	// 关键字类型包括：
	//	- [function]
	//	- IF EQUALS [string] THEN [function] ELSE [function]
	//	- IF NEQ [string] THEN [function] ELSE [function]
	//	- IF GREATER [int] THEN [function] ELSE [function]
	//	- IF LESS [int] THEN [function] ELSE [function]
	//  - FOR [int] TIMES DO [function]
	for _, line := range workflow {
		// 去除行首行尾空格
		line = strings.TrimSpace(line)
		// 如果行为空则跳过
		if line == "" {
			continue
		}
		// 将行按空格分割
		words := strings.Split(line, " ")
		// 根据关键字类型执行对应操作
		switch words[0] {
		case "IF":
			param = handleIf(words, param)
		case "FOR":
			param = handleFor(words, param)
		default:
			param = handleFunction(words[0], param)
		}
	}
}

// handleFunction 处理 function 关键字
func handleFunction(name string, param string) string {
	log.DebugLog("handleFunction: " + name)
	return RunFunction(name, param)
}

// handleIf 处理 IF 关键字
func handleIf(words []string, param string) string {
	log.DebugLog("handleIf")
	switch words[1] {
	case "EQUALS":
		if words[2] == param {
			return handleFunction(words[4], param)
		} else {
			return handleFunction(words[6], param)
		}
	case "NEQ":
		if words[2] != param {
			return handleFunction(words[4], param)
		} else {
			return handleFunction(words[6], param)
		}
	case "GREATER":
		if words[2] > param {
			return handleFunction(words[4], param)
		} else {
			return handleFunction(words[6], param)
		}
	case "LESS":
		if words[2] < param {
			return handleFunction(words[4], param)
		} else {
			return handleFunction(words[6], param)
		}
	default:
		log.ErrorLog("handleIf: unknown keyword")
		return param
	}
}

// handleFor 处理 FOR 关键字
func handleFor(words []string, param string) string {
	log.DebugLog("handleFor")
	// 将循环次数转为整数
	times, err := strconv.Atoi(words[1])
	if err != nil {
		log.ErrorLog("handleFor: " + err.Error())
		return param
	}
	// 循环执行 function
	for i := 0; i < times; i++ {
		param = handleFunction(words[4], param)
	}
	return param
}
