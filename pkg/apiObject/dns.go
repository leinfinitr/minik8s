package apiObject

type Dns struct {
	TypeMeta
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec     DnsSpec    `json:"spec" yaml:"spec"`
	Status  DnsStatus  `json:"status" yaml:"status"`
}
type DnsSpec struct {
	Host string `json:"host" yaml:"host"`
	Paths []Path `json:"paths" yaml:"paths"`
}
type Path struct {
	SubPath string `json:"subPath" yaml:"subPath"`
	SvcIp string `json:"svcIp" yaml:"svcIp"`
	SvcPort string `json:"svcPort" yaml:"svcPort"`
	SvcName string `json:"svcName" yaml:"svcName"`
}
type DnsStatus struct {
	Phase string `json:"phase" yaml:"phase"`
}

type DnsRequest struct {
	Action string `json:"action" yaml:"action"`
	DnsMeta ObjectMeta `json:"dnsMeta" yaml:"dnsMeta"`
}

type HostRequest struct {
	Action string `json:"action" yaml:"action"`
	DnsObject Dns `json:"dnsObject" yaml:"dnsObject"`
	DnsConfig string `json:"dnsConfig" yaml:"dnsConfig"`
	HostList []string `json:"hostList" yaml:"hostList"`
}