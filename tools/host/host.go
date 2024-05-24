package host

import (
	"minik8s/tools/log"
	"net"
	"os"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// GetHostname 获取当前主机的主机名
func GetHostname() (string, error) {
	return os.Hostname()
}

// GetHostIP 获取当前主机的IP
func GetHostIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, i := range interfaces {
		name := i.Name
		if name == "ens33" || name == "eth0" || name == "ens3" {
			// 获取网络接口的地址信息
			addr, err := i.Addrs()
			if err != nil {
				log.ErrorLog("GetHostIP: " + err.Error())
				return "", err
			}

			// 遍历每个地址
			for _, addr := range addr {
				// 检查地址的类型是否为IP地址
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					// 获取IP地址
					return ipNet.IP.String(), nil
				}
			}
		}
	}
	return "", nil
}

// GetTotalMemory 获取当前主机的总内存大小
func GetTotalMemory() (uint64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return memInfo.Total, nil
}

// GetMemoryUsageRate 获取内存使用率
func GetMemoryUsageRate() (float64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return memInfo.UsedPercent, nil
}

// GetCPULoad 获取CPU使用率
func GetCPULoad() ([]float64, error) {
	loads, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}
	return loads, nil
}
