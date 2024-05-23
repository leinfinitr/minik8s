# 在master节点上执行，子网IP为：192.168.1.7
#!/bin/bash
# 清空etcd中的所有数据
etcdctl del "/registry" --prefix
# 重启etcd
systemctl restart etcd
# 启动kube-apiserver
go run minik8s/pkg/apiServer/main
# 启动kubeproxy
go run minik8s/pkg/kubeproxy/main
# 启动kubelet
go run minik8s/pkg/kubelet/main
