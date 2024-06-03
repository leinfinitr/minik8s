package main

import (
	"minik8s/pkg/kubeproxy"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	proxy := kubeproxy.GetKubeproxy()
	proxy.Run()

}
