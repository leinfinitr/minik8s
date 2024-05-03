package etcdclient

import (
	"minik8s/pkg/etcd"
	"minik8s/pkg/klog"
	"minik8s/pkg/config"
)

var EtcdStore *etcd.EtcdClientWrapper = nil

// InitEtcdClient 初始化etcd客户端
func init() {
	etcdConfig := config.NewEtcdConfig()
	etcdclient,err := etcd.NewEtcdClient(etcdConfig.Endpoints,etcdConfig.Timeout)
	if err != nil {
		klog.WarnLog("APIServer","etcd client init failed:"+err.Error())
	}
	EtcdStore = etcdclient
}
