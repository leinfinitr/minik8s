package config

import "minik8s/pkg/apiObject"

const (
	NodesURI      = "/api/v1/nodes"
	NodeURI       = "/api/v1/nodes/:name"
	NodeStatusURI = "/api/v1/nodes/:name/status"

	PodURI                    = "/api/v1/namespaces/:namespace/pods/:name"
	PodEphemeralContainersURI = "/api/v1/namespaces/:namespace/pods/:name/ephemeralcontainers"
	PodLogURI                 = "/api/v1/namespaces/:namespace/pods/:name/log"
	PodStatusURI              = "/api/v1/namespaces/:namespace/pods/:name/status"
	PodsURI                   = "/api/v1/namespaces/:namespace/pods"
	PodsGlobalURI             = "/api/v1/pods"

	ServicesURI      = "/api/v1/namespaces/:namespace/services"
	ServiceURI       = "/api/v1/namespaces/:namespace/services/:name"
	ServiceStatusURI = "/api/v1/namespaces/:namespace/services/:name/status"

	ReplicaSetsURI = "/api/v1/namespaces/:namespace/replicasets"
	GlobalReplicaSetsURI = "/api/v1/replicasets"
	ReplicaSetURI  = "/api/v1/namespaces/:namespace/replicasets/:name"
	ReplicaSetStatusURI = "/api/v1/namespaces/:namespace/replicasets/:name/status"

	HpasURI = "/api/v1/namespaces/:namespace/hpa"
	HpaStatusURI = "/api/v1/namespaces/:namespace/hpa/:name/status"
	HpaURI = "/api/v1/namespaces/:namespace/hpa/:name" 
	GlobalHpaURI = "/api/v1/hpa"
)

const (
	NameSpaceReplace = ":namespace"
	NameReplace      = ":name"
)

var UriMapping = map[string]string{
	apiObject.NodeType: NodesURI,
	apiObject.PodType:  PodsURI,
}

var UriSpecMapping = map[string]string{
	apiObject.NodeType: NodeURI,
	apiObject.PodType:  PodURI,
}