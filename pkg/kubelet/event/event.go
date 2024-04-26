package event

type SyncLoopEventType string

const (
	/* Pod event list
	 * Because pod is smallest unit to deploy, so we concluded event into pod's todo
	 */
	PodNeedStart   SyncLoopEventType = "PodNeedStart"
	PodNeedStop    SyncLoopEventType = "PodNeedStop"
	PodNeedRestart SyncLoopEventType = "PodNeedRestart"
	PodNeedCreate  SyncLoopEventType = "PodNeedCreate"
	PodNeedDelete  SyncLoopEventType = "PodNeedDelete"
	PodNeedUpdate  SyncLoopEventType = "PodNeedUpdate"
)
