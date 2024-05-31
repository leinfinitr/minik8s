package specctlrs

import (
	"minik8s/tools/executor"
	"time"
)

type DnsController interface {
	Run()
}

type DnsControllerImpl struct {
	hostList     []string
	nginxSvcName string
	nginxSvcIp   string
}

func NewDnsController() (DnsController, error) {
	return &DnsControllerImpl{}, nil
}

var(
	DnsControllerDelay   = 0 * time.Second
	DnsControllerTimeGap = []time.Duration{5 * time.Second}
)

func (dc *DnsControllerImpl) Run() {
	// 定期执行
	executor.ExecuteInPeriod(DnsControllerDelay, DnsControllerTimeGap, dc.syncDns)
}

func (dc *DnsControllerImpl) syncDns() {
	// 1. 获取所有的Pod
	// 2. 获取所有的Service
	// 3. 更新DNS
}
