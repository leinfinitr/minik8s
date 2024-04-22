package syncLoop

/* 用于响应并处理收到的SyncLoop时间 */
type SyncLoopHandler struct {
}

func newSyncLoopHandler() *SyncLoopHandler {
	return &SyncLoopHandler{}
}
