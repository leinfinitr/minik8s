package apiObject

type Dns struct {
	TypeMeta
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec     DnsSpec    `json:"spec" yaml:"spec"`
	NginxIP  string     `json:"nginxIP" yaml:"nginxIP"`
}

type DnsSpec struct {
	Host  string `json:"host" yaml:"host"`
	Paths []Path `json:"paths" yaml:"paths"`
}

type Path struct {
	SubPath string `json:"subPath" yaml:"subPath"`
	SvcIp   string `json:"svcIp" yaml:"svcIp"`
	SvcPort string `json:"svcPort" yaml:"svcPort"`
	SvcName string `json:"svcName" yaml:"svcName"`
}

type DnsRequest struct {
	Action  string     `json:"action" yaml:"action"`
	DnsMeta ObjectMeta `json:"dnsMeta" yaml:"dnsMeta"`
}

type HostRequest struct {
	Action    string   `json:"action" yaml:"action"`
	DnsObject Dns      `json:"dnsObject" yaml:"dnsObject"`
	DnsConfig string   `json:"dnsConfig" yaml:"dnsConfig"`
	HostList  []string `json:"hostList" yaml:"hostList"`
}

// Nginx Etcd中存储的数据结构，表示一个用于DNS服务的Nginx pod
type Nginx struct {
	PodIP         string
	Phase         PodPhase
	Name          string
	Namespace     string
	ContainerName string
}
