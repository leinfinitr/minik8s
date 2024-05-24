package main

import (
	"minik8s/pkg/kubeproxy"
)

func main() {
	proxy := kubeproxy.GetKubeproxy()
	proxy.Run()

}
