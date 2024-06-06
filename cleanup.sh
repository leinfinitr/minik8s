#!/bin/bash

# 使用etcdctl删除以/registry为前缀的所有键值
etcdctl del "/registry" --prefix

# 使用crictl删除所有容器（包括正在运行的）
crictl rmp -a -f

# 清空iptables的nat表中的所有规则
iptables -t nat -F

# 删除iptables的nat表中的所有自定义链
iptables -t nat -X
