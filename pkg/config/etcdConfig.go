package config
import "time"

type EtcdConfig struct {
	Endpoints []string
	Timeout   time.Duration
}

const (
	EtcdPodPrefix = "/registry/pods"
	EtcdNodePrefix = "/registry/nodes"
	EtcdReplicaSetPrefix = "/registry/replicasets"
	EtcdHpaPrefix = "/registry/hpa"
)
func NewEtcdConfig() *EtcdConfig {
	return &EtcdConfig{
		Endpoints: []string{"localhost:2379"},
		Timeout:   3 * time.Second,
	}
}
