package main

import "minik8s/pkg/controller/ctlrMgr"

func main() {
	ctrlManager := ctlrmgr.NewControllerManager()

	ctrlManager.Run(make(chan struct{}))
}
