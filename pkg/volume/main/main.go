package main

import (
	"fmt"
	"minik8s/pkg/volume/storageManager"
)

func main() {
	manager := storageManager.NewStorageManager("/mnt/local-storage")
	err := manager.CreatePV("pv1", 1024*1024*1024, []string{"ReadWriteOnce"})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = manager.CreatePVC("pvc1", "default", 512*1024*1024, []string{"ReadWriteOnce"})
	if err != nil {
		fmt.Println(err)
		return
	}
	// 打印结果检查
	fmt.Println(manager.PVs)
	fmt.Println(manager.PVCs)
}
