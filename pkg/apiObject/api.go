// 描述: 定义api对象的基本结构
// 参考：https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/apis/meta/v1/types.go

package apiObject

type TypeMeta struct {
	// Kind: 对象的类型，如Pod、Service、ReplicationController
	Kind string `json:"kind" yaml:"kind"`
	// APIVersion: API版本
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
}

type ObjectMeta struct {
	// Name: 对象的名称
	Name string `json:"name" yaml:"name"`
	// Namespace: 对象所在的命名空间
	Namespace string `json:"namespace" yaml:"namespace"`
	// Labels: 对象的标签
	Labels map[string]string `json:"labels" yaml:"labels"`
}
