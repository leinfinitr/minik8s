#!/bin/bash
# 清空etcd中的所有数据
etcdctl del "/registry" --prefix
# 重启etcd
systemctl restart etcd
