// 描述: Service对象的封装
// 参考：https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/api/core/v1/types.go

package apiObject

type Service struct {
	// 对象的类型元数据
	TypeMeta
	// 对象的元数据
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	// Service的规格
	Spec ServiceSpec `json:"spec" yaml:"spec"`
	// Service的状态
	Status ServiceStatus `json:"status" yaml:"status"`
}

type ServiceSpec struct {
	// 该Service暴露的端口列表，可以对一个服务指定多个端口
	Ports []ServicePort `json:"ports" yaml:"ports"`
	// 将service流量路由到具有与此selector匹配的标签键和值的pod。
	//	如果Type为ExternalName，则忽略此字段
	Selector map[string]string `json:"selector" yaml:"selector"`
	// Service的IP地址，通常是随机分配的
	ClusterIP string `json:"clusterIP" yaml:"clusterIP"`
	// Service的类型，包括ClusterIP、NodePort、LoadBalancer、ExternalName，默认为ClusterIP
	Type string `json:"type" yaml:"type"`
}

type ServicePort struct {
	// Service端口的名称
	Name string `json:"name" yaml:"name"`
	// 该端口的IP协议，包括TCP、UDP、SCTP，默认为TCP
	Protocol Protocol `json:"protocol" yaml:"protocol"`
	// Service内部向外暴露的端口
	Port int32 `json:"port" yaml:"port"`
	// 服务所针对的pod上要访问的端口的编号或名称，也就是对应到pod上的端口，如果未指定，则使用port作为targetPort
	TargetPort int32 `json:"targetPort" yaml:"targetPort"`
	// 当Type为NodePort或LoadBalancer时，每个Node上的端口，即全局对外提供服务的端口
	NodePort int32 `json:"nodePort" yaml:"nodePort"`
}

type ServiceStatus struct {
	// LoadBalancer的状态
	LoadBalancer LoadBalancerStatus `json:"loadBalancer" yaml:"loadBalancer"`
	// 当前Service的状态
	Conditions []Condition
}

type LoadBalancerStatus struct {
	// Ingress是一个包含所有LoadBalancer的地址的列表
	Ingress []LoadBalancerIngress `json:"ingress" yaml:"ingress"`
}

type LoadBalancerIngress struct {
	// LoadBalancer的IP地址
	IP string `json:"ip" yaml:"ip"`
	// LoadBalancer的主机名
	Hostname string `json:"hostname" yaml:"hostname"`
}

type Condition struct {
	// Condition的类型
	// 	Ready: Service是否已准备好
	// 	NetworkUnavailable: Service是否无法访问网络
	Type string `json:"type" yaml:"type"`
	// Condition的状态
	// 	True: 条件满足
	// 	False: 条件不满足
	// 	Unknown: 状态未知
	Status string `json:"status" yaml:"status"`
}

// Endpoint pod信息的子集，用来实现service和pod的映射，给Kubeproxy使用
type Endpoint struct {
	PodUUID string
	IP      string
}
