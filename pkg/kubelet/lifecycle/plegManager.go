package lifecycle

// PLEG(Pod Lifecycle Event Generator）官方参考功能：
// 	定时调用 container runtime 获取本节点 containers/sandboxes 的信息，
//	并与自身维护的 pods cache 信息进行对比，生成对应的 PodLifecycleEvent，
//	然后输出到 eventChannel 中，通过 eventChannel 发送到 kubelet syncLoop 进行消费，
//	然后由 kubelet syncPod 来触发 pod 同步处理过程，最终达到用户的期望状态。

// PlegManager
// 1. 负责pod生命周期的管理
type PlegManager interface {
}
