package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/pkg/kubectl/translator"
	httprequest "minik8s/tools/httpRequest"
	"minik8s/tools/log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var serverlessCmd = &cobra.Command{
	Use:   "serverless",
	Short: "Create, delete, and get serverless functions",
	Long:  "Create, delete, and get serverless functions",
	Run:   serverlessHandler,
}

func serverlessHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		// 如果没有参数，输出 serverless 指令的帮助信息
		printHelp()
		return
	}

	// 如果有参数，根据参数执行相应的操作
	switch args[0] {
	case "create":
		if len(args) != 2 {
			log.ErrorLog("The number of parameters is incorrect: create [file.yaml]")
			return
		}
		createServerless(args[1])
	case "delete":
		if len(args) != 2 {
			log.ErrorLog("The number of parameters is incorrect: delete [serverless name]")
			return
		}
		deleteServerless(args[1])
	case "get":
		if len(args) != 1 {
			log.ErrorLog("The number of parameters is incorrect: get")
			return
		}
		getAllServerless()
	case "update":
		if len(args) != 3 {
			log.ErrorLog("The number of parameters is incorrect: update [serverless name] [file.py]")
			return
		}
		updateFunction(args[1], args[2])
	case "run":
		if len(args) != 3 {
			log.ErrorLog("The number of parameters is incorrect: run [serverless name] [param]")
			return
		}
		runFunction(args[1], args[2])
	case "workflow":
		if len(args) != 2 {
			log.ErrorLog("The number of parameters is incorrect: workflow [file.txt] [param]")
			return
		}
		workflow(args[1], args[2])
	default:
		printHelp()
	}
}

// printHelp 输出 serverless 指令的帮助信息
func printHelp() {
	help := `serverless 指令格式如下：
	serverless create [xxx.yaml] # 创建 serverless 环境
	serverless delete [serverless name] # 删除 serverless 环境
	serverless get # 获取所有 serverless 环境
	serverless update [serverless name] [xxx.py] # 更新指定 serverless 环境函数
	serverless run [serverless name] [param] # 运行指定 serverless 环境
	serverless workflow [xxx.txt] [param] # 执行 serverless 工作流`
	println(help)
}

// createServerless 创建 serverless 环境
func createServerless(fileName string) {
	log.DebugLog("Create serverless environment: " + fileName)
	// 检查文件是否存在
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		log.ErrorLog("The resource type specified does not exist.")
		os.Exit(1)
	}
	if fileInfo.IsDir() {
		log.ErrorLog("The resource type specified is a directory.")
		os.Exit(1)
	}
	// 读取文件内容
	content, err := os.ReadFile(fileName)
	if err != nil {
		log.ErrorLog("Could not read the file specified.")
		os.Exit(1)
	}
	// 解析 yaml 文件
	var serverless apiObject.Serverless
	err = translator.ParseApiObjFromYaml(content, &serverless)
	if err != nil {
		log.ErrorLog("Could not unmarshal the yaml file.")
		os.Exit(1)
	}
	// 转发给 serverless 服务端口处理
	url := config.ServerlessURL() + config.ServerlessURI
	res, err := httprequest.PostObjMsg(url, serverless)
	if err != nil {
		log.ErrorLog("Could not post the object message." + err.Error())
		os.Exit(1)
	}
	// 输出结果
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.ErrorLog("Could not read the response body.")
		os.Exit(1)
	}
	fmt.Println("Result: " + string(body))
}

// getAllServerless 获取所有 serverless 环境
func getAllServerless() {
	log.DebugLog("Get all serverless environment")
	var result []apiObject.Serverless
	// 转发给 serverless 服务端口处理
	url := config.ServerlessURL() + config.ServerlessURI
	response, err := httprequest.GetMsg(url)
	if err != nil {
		log.ErrorLog(err.Error())
		os.Exit(1)
	}
	// 解析返回结果
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.ErrorLog(err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.ErrorLog(err.Error())
		os.Exit(1)
	}
	// 输出结果
	for _, serverless := range result {
		fmt.Println("Name: " + serverless.Name + " Image: " + serverless.Image + " Volume: " + serverless.Volume)
	}
}

// deleteServerless 删除 serverless function
func deleteServerless(serverlessName string) {
	log.DebugLog("Delete serverless environment: " + serverlessName)
	// 转发给 serverless 服务端口处理
	url := config.ServerlessURL() + config.ServerlessFunctionURI
	url = strings.Replace(url, config.NameReplace, serverlessName, -1)
	res, err := httprequest.DelMsg(url)
	if err != nil {
		log.ErrorLog(err.Error())
		os.Exit(1)
	}
	// 输出结果
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.ErrorLog(err.Error())
		os.Exit(1)
	}
	fmt.Println("Result: " + string(body))
}

// updateFunction 更新 serverless 环境函数
func updateFunction(serverlessName, fileName string) {
	log.DebugLog("Update serverless: " + serverlessName + " " + fileName)
	// 检查文件是否存在
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		log.ErrorLog("The file " + fileName + " does not exist.")
		os.Exit(1)
	}
	if fileInfo.IsDir() {
		log.ErrorLog("The file " + fileName + " is a directory.")
		os.Exit(1)
	}
	// 转发给 serverless 服务端口处理
	url := config.ServerlessURL() + config.ServerlessFunctionURI
	url = strings.Replace(url, config.NameReplace, serverlessName, -1)
	response, err := httprequest.PutObjMsg(url, fileName)
	if err != nil {
		log.ErrorLog(err.Error())
		os.Exit(1)
	}
	// 输出结果
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.ErrorLog(err.Error())
		os.Exit(1)
	}
	fmt.Println("Result: " + string(body))
}

// runFunction 运行 serverless 环境函数
func runFunction(serverlessName string, param string) {
	log.DebugLog("Run serverless: " + serverlessName)
	// 转发给 serverless 服务端口处理
	url := config.ServerlessURL() + config.ServerlessRunURI
	url = strings.Replace(url, config.NameReplace, serverlessName, -1)
	url = strings.Replace(url, config.ParamReplace, param, -1)
	response, err := httprequest.GetMsg(url)
	if err != nil {
		log.ErrorLog(err.Error())
		os.Exit(1)
	}
	// 输出结果
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.ErrorLog(err.Error())
		os.Exit(1)
	}
	fmt.Println("Result: " + string(body))
}

// workflow 执行 serverless 工作流
func workflow(fileName string, param string) {
	log.DebugLog("Run workflow: " + fileName)
	// 检查文件是否存在
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		log.ErrorLog("The file " + fileName + " does not exist.")
		os.Exit(1)
	}
	if fileInfo.IsDir() {
		log.ErrorLog("The file " + fileName + " is a directory.")
		os.Exit(1)
	}
	// 读取文件内容
	content, err := os.ReadFile(fileName)
	if err != nil {
		log.ErrorLog("Could not read the file specified.")
		os.Exit(1)
	}
	// 转发给 serverless 服务端口处理
	url := config.ServerlessURL() + config.ServerlessWorkflowURI
	url = strings.Replace(url, config.ParamReplace, param, -1)
	response, err := httprequest.PostObjMsg(url, content)
	if err != nil {
		log.ErrorLog(err.Error())
		os.Exit(1)
	}
	// 输出结果
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.ErrorLog(err.Error())
		os.Exit(1)
	}
	fmt.Println("Result: " + string(body))
}