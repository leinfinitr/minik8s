package config

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
)

const (
	NameSpaceReplace = ":namespace"
	NameReplace      = ":name"
)
