package apiObject

type Serverless struct {
	// 所需的python镜像
	Image string `json:"image" yaml:"image"`
	// 名称
	Name string `json:"name" yaml:"name"`
	// 存放函数文件的 volume
	Volume string `json:"volume" yaml:"volume"`
	// 函数文件所在主机目录
	HostPath string `json:"hostPath" yaml:"hostPath"`
	// 函数文件名称
	FunctionFile string `json:"functionFile" yaml:"functionFile"`
	// 所需环境
	Requirements string `json:"requirements" yaml:"requirements"`
	// 执行函数的命令
	Command string `json:"command" yaml:"command"`
}
