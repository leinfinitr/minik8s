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
	PodsURI                   = "/api/v1/namespaces/:namespace/pods"
	PodsGlobalURI             = "/api/v1/pods"

	ProxyStatusURI   = "/api/v1/proxy"
	ProxiesStatusURI = "/api/v1/proxy/:name"

	ServicesURI      = "/api/v1/namespaces/:namespace/services"
	ServiceURI       = "/api/v1/namespaces/:namespace/services/:name"
	ServiceStatusURI = "/api/v1/namespaces/:namespace/services/:name/status"

	ServerlessURI         = "/api/v1/serverless"
	ServerlessFunctionURI = "/api/v1/serverless/function/:name"
	ServerlessRunURI      = "/api/v1/serverless/run/:name/:param"
	ServerlessWorkflowURI = "/api/v1/serverless/workflow/:param"
)

const (
	NameSpaceReplace = ":namespace"
	NameReplace      = ":name"
	ParamReplace     = ":param"
)

var UriMapping = map[string]string{
	apiObject.NodeType: NodesURI,
	apiObject.PodType:  PodsURI,
}

var UriSpecMapping = map[string]string{
	apiObject.NodeType: NodeURI,
	apiObject.PodType:  PodURI,
}
