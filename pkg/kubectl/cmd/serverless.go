package cmd

import (
	"github.com/spf13/cobra"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/pkg/kubectl/translator"
	"minik8s/tools/httpRequest"
	"minik8s/tools/log"
	"os"
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
			log.ErrorLog("The number of parameters is incorrect.")
			return
		}
		createServerless(args[1])
	//case "delete":
	//	deleteServerless(args[1])
	//case "get":
	//	getServerless(args[1])
	//case "add":
	//	addFunction(args[1], args[2])
	//case "remove":
	//	removeFunction(args[1], args[2])
	//case "update":
	//	updateFunction(args[1], args[2])
	default:
		printHelp()
	}
}

// printHelp 输出 serverless 指令的帮助信息
func printHelp() {
	help := `serverless 指令格式如下：
	serverless create xxx.yaml # 创建 serverless 环境
	serverless delete [serverless name] # 删除 serverless 环境
	serverless get # 获取所有 serverless 环境
	serverless get [serverless name] # 获取指定 serverless 环境
	serverless add [serverless name] xxx.py # 向指定 serverless 环境添加函数
	serverless remove [serverless name] [function name] # 从指定 serverless 环境删除函数
	serverless update [serverless name] xxx.py # 向指定 serverless 环境添加函数`
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
	_, err = httprequest.PostObjMsg(url, serverless)
	if err != nil {
		log.ErrorLog("Could not post the object message." + err.Error())
		os.Exit(1)
	}
}
