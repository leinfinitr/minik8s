package config
import "strconv"

const (
	SchedulerConfigPath = "/scheduler"
	SchedulerLocalAddress = "127.0.0.1"
	SchedulerLocalPort = 7820
)

func SchedulerURL() string {
	return "http://" + SchedulerLocalAddress + ":" + strconv.Itoa(SchedulerLocalPort)
}
func SchedulerPort() string{
	return strconv.Itoa(SchedulerLocalPort)
}

func SchedulerPath() string {
	return SchedulerConfigPath
}
func NewSchedulerConfig() *SchedulerConfig {
	return &SchedulerConfig{
		SchedulerIP: SchedulerLocalAddress,
		SchedulerPort: SchedulerLocalPort,
	}
}

type SchedulerConfig struct {
	SchedulerIP string
	SchedulerPort int
}

func (c *SchedulerConfig) SchedulerURL() string {
	return "http://" + c.SchedulerIP + ":" + strconv.Itoa(c.SchedulerPort)
}