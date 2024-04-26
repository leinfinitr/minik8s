package lifecycle

/* pleg：Pod-Lifecycle-Event-Generator，其功能如下
 * 1. 负责pod生命周期的关系
 * 2. 会定时获得本节点的运行信息，并与维护的相关信息相对比
 * 3. 根据对比情况生成LifeCycleEvent，放入队列中，由syncLooper处理
 */
type PlegManager interface {
}
