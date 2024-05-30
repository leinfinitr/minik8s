package apiObject

const (
	ServerlessEventTypeTime ServerlessEventType = "time"
	ServerlessEventTypeFile ServerlessEventType = "file"
)

type Serverless struct {
	// 所需的python镜像
	Image string `json:"image" yaml:"image"`
	// 名称
	Name string `json:"name" yaml:"name"`
	// 函数文件所在主机目录
	HostPath string `json:"hostPath" yaml:"hostPath"`
	// 函数文件名称
	FunctionFile string `json:"functionFile" yaml:"functionFile"`
	// 所需环境
	Requirements string `json:"requirements" yaml:"requirements"`
	// 执行函数的命令
	Command string `json:"command" yaml:"command"`
}

type ServerlessEvent struct {
	// 事件类型
	Type ServerlessEventType `json:"type" yaml:"type"`
	// 事件绑定参数
	Params string `json:"params" yaml:"params"`
	// 事件触发的 Serverless Function 名称
	Name string `json:"name" yaml:"name"`
	// 事件触发的 Serverless Function 参数
	Args string `json:"args" yaml:"args"`
}

type ServerlessEventType string
