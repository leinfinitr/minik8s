// 测试host相关的操作

package host

import (
	"github.com/stretchr/testify/assert"
	"minik8s/tools/log"
	"strconv"
	"testing"
)

func TestGetHostname(t *testing.T) {
	hostname, err := GetHostname()
	log.InfoLog(hostname)
	assert.NoError(t, err)
	assert.NotEmpty(t, hostname)
}

func TestGetIP(t *testing.T) {
	ips, err := GetHostIP()
	log.InfoLog(ips)
	assert.NoError(t, err)
}

func TestGetTotalMemory(t *testing.T) {
	totalMemory, err := GetTotalMemory()
	log.InfoLog(strconv.FormatUint(totalMemory, 10))
	assert.NoError(t, err)
	assert.True(t, totalMemory > 0)
}

func TestGetMemoryUsageRate(t *testing.T) {
	memoryUsageRate, err := GetMemoryUsageRate()
	log.InfoLog(strconv.FormatFloat(memoryUsageRate, 'f', -1, 64))
	assert.NoError(t, err)
	assert.True(t, memoryUsageRate >= 0 && memoryUsageRate <= 100)
}

func TestGetCPULoad(t *testing.T) {
	cpuLoad, err := GetCPULoad()
	log.InfoLog(strconv.FormatFloat(cpuLoad[0], 'f', -1, 64))
	assert.NoError(t, err)
	assert.NotEmpty(t, cpuLoad)
}
