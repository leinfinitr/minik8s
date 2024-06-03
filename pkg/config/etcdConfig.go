package config

import "time"

type EtcdConfig struct {
	Endpoints []string
	Timeout   time.Duration
}

const (
	EtcdPodPrefix        = "/registry/pods"
	EtcdNodePrefix       = "/registry/nodes"
	EtcdReplicaSetPrefix = "/registry/replicasets"
	EtcdHpaPrefix        = "/registry/hpa"
	EtcdServerlessPrefix = "/registry/serverless"
	EtcdServicePrefix    = "/registry/services"
	EtcdJobPrefix        = "/registry/jobs"
	EtcdJobCodePrefix    = "/registry/jobCode"
	EtcdPvPrefix         = "/registry/pv"
	EtcdPvcPrefix        = "/registry/pvc"
	EtcdDnsPrefix        = "/registry/dns"
	EtcdDnsRequestPrefix = "/registry/dnsrequest"
)

func NewEtcdConfig() *EtcdConfig {
	return &EtcdConfig{
		Endpoints: []string{"localhost:2379"},
		Timeout:   3 * time.Second,
	}
}
