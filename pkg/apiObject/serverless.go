package apiObject

type Serverless struct {
	// 所需的python镜像
	Image string `json:"image" yaml:"image"`
	// 名称
	Name string `json:"name" yaml:"name"`
	// 函数
	Function []ServerlessFunction `json:"function" yaml:"function"`
}

// 创建函数时指定的方法
const (
	CreateFunctionByCode = "code"
	CreateFunctionByFile = "file"
)

type ServerlessFunction struct {
	// 函数名称
	Name string `json:"name" yaml:"name"`
	// 创建函数的方法
	CreateMethod string `json:"createMethod" yaml:"createMethod"`
	// 函数代码
	Code string `json:"code" yaml:"code"`
	// 函数文件
	File string `json:"file" yaml:"file"`
}
