package apiObject

type Serverless struct {
	// 所需的python镜像
	Image string `json:"image" yaml:"image"`
	// 名称
	Name string `json:"name" yaml:"name"`
	// 存放函数文件的 volume
	Volume string `json:"volume" yaml:"volume"`
	// 所需环境
	Requirements []string `json:"requirements" yaml:"requirements"`
	// 函数文件
	FunctionFile string `json:"functionFile" yaml:"functionFile"`
}
