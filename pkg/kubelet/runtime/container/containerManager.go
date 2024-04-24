package container

/* containerManager 负责 node 节点上运行的容器的 cgroup 配置信息，
 * kubelet 启动参数如果指定 --cgroups-per-qos 的时候，kubelet 会启动
 * goroutine 来周期性的更新 pod 的 cgroup 信息，维护其正确性，该参数默认为 true,
 * 实现了 pod 的Guaranteed/BestEffort/Burstable 三种级别的 Qos。
 */
type ContainerManager struct {
}
