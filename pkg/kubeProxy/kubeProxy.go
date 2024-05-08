// 描述：kubeProxy 用于创建一个代理服务器，用于转发请求到后端服务
// 	1. 接受一个pod的虚拟IP，检查其是否属于该节点

package kubeProxy

type KubeProxy struct {
	// IpTables 存储每个Node的IP地址及其对应Pod的虚拟IP地址段
	IpTables map[string][]string
}
