// 描述: 定义api对象的基本结构
// 参考：https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/apis/meta/v1/types.go

package apiObject

type TypeMeta struct {
	// 对象的类型，如Pod、Service、ReplicationController
	Kind string `json:"kind" yaml:"kind"`
	// API版本
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
}

type ObjectMeta struct {
	// 对象的名称
	Name string `json:"name" yaml:"name"`
	// 对象所在的命名空间
	Namespace string `json:"namespace" yaml:"namespace"`
	// 对象的标签
	Labels map[string]string `json:"labels" yaml:"labels"`
	// 对象的注解
	Annotations map[string]string `json:"annotations" yaml:"annotations"`
	// 此对象在时间和空间上唯一的值。它通常由成功创建资源时的服务器生成，并且不允许在PUT上的更改操作。
	UUID string `json:"uid" yaml:"uid"`
}

const (
	PodType         = "Pod"
	ServiceType     = "Service"
	ReplicaSetType  = "ReplicaSet"
	NodeType        = "Node"
	HpaType		 	= "Hpa"
	ContainerType  = "Container"
)

var AllTypeList = []string{PodType, ServiceType, ReplicaSetType, NodeType}