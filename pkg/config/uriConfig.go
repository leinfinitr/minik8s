package config

import "minik8s/pkg/apiObject"

const (
	NodesURI      = "/api/v1/nodes"
	NodeURI       = "/api/v1/nodes/:name"
	NodeStatusURI = "/api/v1/nodes/:name/status"

	PodURI                    = "/api/v1/namespaces/:namespace/pods/:name"
	PodEphemeralContainersURI = "/api/v1/namespaces/:namespace/pods/:name/ephemeralContainers"
	PodLogURI                 = "/api/v1/namespaces/:namespace/pods/:name/log"
	PodStatusURI              = "/api/v1/namespaces/:namespace/pods/:name/status"
	PodExecURI                = "/api/v1/namespaces/:namespace/pods/:name/exec/:container/:param"
	PodsURI                   = "/api/v1/namespaces/:namespace/pods"
	PodsGlobalURI             = "/api/v1/pods"
	PodsSyncURI               = "/api/v1/pods/sync"

	ProxyStatusURI   = "/api/v1/proxy"
	ProxiesStatusURI = "/api/v1/proxy/:name"

	ServicesURI      = "/api/v1/namespaces/:namespace/services"
	ServiceURI       = "/api/v1/namespaces/:namespace/services/:name"
	ServiceStatusURI = "/api/v1/namespaces/:namespace/services/:name/status"

	ReplicaSetsURI       = "/api/v1/namespaces/:namespace/replicasets"
	GlobalReplicaSetsURI = "/api/v1/replicasets"
	ReplicaSetURI        = "/api/v1/namespaces/:namespace/replicasets/:name"
	ReplicaSetStatusURI  = "/api/v1/namespaces/:namespace/replicasets/:name/status"

	HpasURI      = "/api/v1/namespaces/:namespace/hpa"
	HpaStatusURI = "/api/v1/namespaces/:namespace/hpa/:name/status"
	HpaURI       = "/api/v1/namespaces/:namespace/hpa/:name"
	GlobalHpaURI = "/api/v1/hpa"

	ServerlessURI         = "/api/v1/serverless"
	ServerlessFunctionURI = "/api/v1/serverless/function/:name"
	ServerlessRunURI      = "/api/v1/serverless/run/:name/:param"
	ServerlessWorkflowURI = "/api/v1/serverless/workflow/:param"

	JobsURI      = "/api/v1/namespaces/:namespace/jobs"
	JobURI       = "/api/v1/namespaces/:namespace/jobs/:name"
	GlobalJobURI = "/api/v1/jobs"
	JobStatusURI = "/api/v1/namespaces/:namespace/jobs/:name/status"
	JobsCodeURI  = "/api/v1/namespaces/:namespace/jobs/code"
	JobCodeURI   = "/api/v1/namespaces/:namespace/jobs/:name/code"

	PersistentVolumeURI  = "/api/v1/pv"
	PersistentVolumesURI = "/api/v1/pv/:name"

	PersistentVolumeClaimURI  = "/api/v1/pvc"
	PersistentVolumeClaimsURI = "/api/v1/pvc/:name"

	MonitorURL = "/api/v1/monitor"
)

const (
	NameSpaceReplace = ":namespace"
	NameReplace      = ":name"
	ParamReplace     = ":param"
	ContainerReplace = ":container"
)

var UriMapping = map[string]string{
	apiObject.NodeType: NodesURI,
	apiObject.PodType:  PodsURI,
}

var UriSpecMapping = map[string]string{
	apiObject.NodeType: NodeURI,
	apiObject.PodType:  PodURI,
}
