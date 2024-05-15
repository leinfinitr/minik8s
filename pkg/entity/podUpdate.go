package entity

type PodUpdateCmd struct {
	Action string
	// PodTarget apiObject.PodStore
	Node string
	Cmd  []string
}
