package main

import "minik8s/pkg/controller"

func main() {
	ctrlManager := controller.NewControllerManager()

	ctrlManager.Run(make(chan struct{}))
}
