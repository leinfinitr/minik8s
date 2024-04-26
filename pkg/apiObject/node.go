package apiObject

type Node struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	IP         string `json:"ip" yaml:"ip"`
	NodeMeta   NodeMetaData
}

// 专门给用户呈现的Node数据结构，在这里不包含namespace
type NodeMetaData struct {
	UUID        string            `json:"uuid" yaml:"uuid"`
	Name        string            `json:"name" yaml:"name"`
	Labels      map[string]string `json:"labels" yaml:"labels"`
	Annotations map[string]string `json:"annotations" yaml:"annotations"`
}
