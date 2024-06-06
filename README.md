# MiniK8s 云操作系统
**MiniK8s**是以Kubernetes为原型的简化版云操作系统，能够提供负载均衡、流量控制、弹性伸缩、容错等功能。

## 结构设计
项目整体架构如下图所示
<div align="middle">
<img src=./docs/assets/minik8s.png width=100%>
</div>

### Master Node
- **API Server**：minik8s 与客户端和其他组件交互的核心组件。
- **Etcd**：存储 Pod、Container 等资源的元数据，负责持久化存储。
- **Scheduler**：接受来自 API Server 的调度请求，采用轮询的调度策略得到目标节点并返回。
- **Controller Manager**：许多功能组件的集合体。
  - **HPA Controller**：监控cpu和memory的资源占用，并根据负载高低调整副本数量。
  - **ReplicaSet Controller**：实现ReplicaSet的资源实现。
  - **Job Controller**：监控并维护Job对象到实际资源对象的实现
  - **PV Controller**：用于接受来自 API Server 的创建持久化卷声明、持久化卷的绑定和解绑等功能。
- **Serverless**：提供 Serverless 的创建、运行、扩缩容控制、控制流分析执行等所有与 Serverless 相关的功能。
### Worker Node
- **Kubeproxy**：负责各个节点的网络配置，提供负载均衡、流量转发等功能。
  - **Iptable Manager**：动态更新iptables，维护service 到Pod IP的映射。
- **Kubelet**：minik8s在各个Worker Node的中心组件，负责Pod等资源的生命周期管理。
  - **Runtime**：通过CRI接口去和containerd进行交互，包括创建Pod、container、更新状态等。
  - **Image Manager**：负责对容器所需镜像的管理。
  - **Pod Manager**：负责对Pod资源生命周期的管理
### 其他服务工具
- **Kubectl**：minik8s 的客户端工具，用于接收和分析用户指令，进行格式检查和筛选，进而转发给 API Server 或 Serverless。
- **Monitor**：结合nodeExporter、Grafana、Prometheus等组件实现对节点和pod的资源的监控，并且提供炫酷的可视化界面。
- **Nginx Pod**：提供在系统的DNS服务，实现域名下的子路径到Service的流量转发。
- **GPU-ServerPod**：提供远程调用GPU运行程序的支持。


## 文档仓库
主要存放我们在开发过程中的相关文档，相关的飞书链接我们已经公开。

### 功能文档

- 基础功能
  - [Pod抽象](https://kxd3r8u0zxd.feishu.cn/wiki/I8a3w8EXGiFD3YkVaL7cHa9SnMf?from=from_copylink)：包括Pod等对象的抽象实现。
  - [CRI实现](https://kxd3r8u0zxd.feishu.cn/wiki/YxoNwy0yei9CCaktZFdc9AlgnRg?from=from_copylink)：底层与containerd交互的相关实现。
  - [Service](https://kxd3r8u0zxd.feishu.cn/wiki/Q3wZwrQdoikn45kLkRzcaoJBneb?from=from_copylink): Service的设计与支持。
  - [DNS转发](https://kxd3r8u0zxd.feishu.cn/wiki/NuwywPzB2iFpXwkmYSNc5qGSnXl?from=from_copylink): DNS功能的设计与支持。
  - [ReplicaSet](https://kxd3r8u0zxd.feishu.cn/wiki/UAiJwcyMyi5eAyk9Oydcr2GGnYf?from=from_copylink): Workload功能的设计与支持。
  - [HPA](https://kxd3r8u0zxd.feishu.cn/wiki/RjomwaFxdiMscYkLX0ScFBNIngc?from=from_copylink)：HPA功能的设计与实现。
- 进阶功能
  - [监控与日志](https://kxd3r8u0zxd.feishu.cn/wiki/ZI3Sw3pExi4smmkpiYVcTjvlntc?from=from_copylink)：通过各种组件来提供对系统的监控与日志功能。
  - [容错](https://kxd3r8u0zxd.feishu.cn/wiki/Tb7Kwdw3JiPFt1kxnBFcElIKnEd?from=from_copylink)：分析了可能的各种容错情况以及实现。
  - [PV&PVC](https://kxd3r8u0zxd.feishu.cn/wiki/E0X5wj2lriDJMZk2ESwcXTkenSf?from=from_copylink)：关于持久化卷和持久化声明的实现。

### 其他文档
- [上机测试](https://kxd3r8u0zxd.feishu.cn/wiki/L1MTwP7emixvMNkIwhBcZYKfnef?from=from_copylink)：如果在机器上启动minik8s，以及相关测试流程。
- [环境配置](https://kxd3r8u0zxd.feishu.cn/wiki/Q7dWwAF1BiFruhkTVCucv7QFn0c?from=from_copylink)：在启动minik8s前，需要对环境进行配置。
- [CNI网络配置](https://kxd3r8u0zxd.feishu.cn/wiki/EKRQwtvUQiPkUAkFI9RcTw9PnPc?from=from_copylink)：此处我们采用CNI 插件Flannel来实现全局的网络配置。


## 软件栈
| 软件栈                  | 提供功能                           |
|-------------------------|------------------------------------|
| google/uuid             | 为 pod 等结构体的元数据生成 uuid   |
| spf13/cobra             | 对客户端输入的命令进行解析         |
| distribution/reference  | 为镜像生成 Tag                     |
| coreos/go-iptables      | 用于 kubeproxy 对于 iptables 的管理 |
| shirou/gopsutil         | 获取主机内存、CPU资源和使用率       |
| prometheus/client_golang| 提供 prometheus 监控支持           |
| fsnotify/fsnotify       | Serverless 中对文件的监控          |
| jedib0t/go-pretty       | 美化命令行的输出                   |
| stretchr/testify        | 用于测试中的断言支持               |
| fatih/color             | 美化 log 输出                      |

## 项目贡献
|   姓名  | 身份 |  贡献度 	| 功能实现 |
|------	|------|----------|----------------------------|
| 刘佳隆 | 组长  |   33.3% 	|Pod抽象、Serverless、PV与PVC等|
| 吴银朋 | 组员  |   33.3%	|Service、DNS转发、容错、日志与监控等|
| 顾泽钜 | 组员  |   33.3%	|Scheduler、ReplicaSet、HPA、GPU等 |