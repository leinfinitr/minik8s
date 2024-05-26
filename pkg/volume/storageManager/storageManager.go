package storageManager

import (
	"fmt"
	"minik8s/pkg/apiObject"
	"os"
	"path/filepath"
)

// StorageManager 存储管理器
type StorageManager struct {
	PVs      map[string]*apiObject.PersistentVolume
	PVCs     map[string]*apiObject.PersistentVolumeClaim
	BasePath string // 本地文件系统的基础路径
}

func NewStorageManager(basePath string) *StorageManager {
	return &StorageManager{
		PVs:      make(map[string]*apiObject.PersistentVolume),
		PVCs:     make(map[string]*apiObject.PersistentVolumeClaim),
		BasePath: basePath,
	}
}

// CreatePV 静态创建PV
func (sm *StorageManager) CreatePV(name string, capacity int64, accessModes []string) error {
	pv := &apiObject.PersistentVolume{
		Name:        name,
		Capacity:    capacity,
		Path:        filepath.Join(sm.BasePath, name),
		AccessModes: accessModes,
		IsBound:     false,
		ClaimedBy:   "",
	}
	if err := os.MkdirAll(pv.Path, 0755); err != nil {
		return err
	}
	sm.PVs[name] = pv
	return nil
}

// CreatePVC 创建PVC
func (sm *StorageManager) CreatePVC(name, namespace string, capacity int64, accessModes []string) error {
	pvc := &apiObject.PersistentVolumeClaim{
		Name:        name,
		Namespace:   namespace,
		Capacity:    capacity,
		AccessModes: accessModes,
		BoundTo:     "",
	}
	sm.PVCs[name] = pvc
	return sm.bindPVCIfPossible(pvc)
}

// bindPVCIfPossible 尝试绑定PVC到合适的PV
func (sm *StorageManager) bindPVCIfPossible(pvc *apiObject.PersistentVolumeClaim) error {
	// 简单的逻辑，根据容量和访问模式匹配PV
	for _, pv := range sm.PVs {
		if !pv.IsBound && pv.Capacity >= pvc.Capacity &&
			equalAccessModes(pv.AccessModes, pvc.AccessModes) {
			pv.IsBound = true
			pv.ClaimedBy = pvc.Name
			pvc.BoundTo = pv.Name
			return nil
		}
	}
	return fmt.Errorf("no suitable PV found for PVC %s", pvc.Name)
}

// equalAccessModes 简单比较两个访问模式列表是否相等
func equalAccessModes(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
