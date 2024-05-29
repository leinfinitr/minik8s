package specctlrs

import (
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/executor"
	"minik8s/tools/log"
	"os/exec"
	"time"
)

type PvController interface {
	Run()
	AddPv(pv *apiObject.PersistentVolume) error
	AddPvc(pvc *apiObject.PersistentVolumeClaim) error
	BindPvc(pvc *apiObject.PersistentVolumeClaim) error
}

type PvControllerImpl struct {
}

var (
	PvControllerDelay   = 3 * time.Second
	PvControllerTimeGap = []time.Duration{10 * time.Second}
)

var PvControllerInstance *PvControllerImpl = nil

func NewPvController() (PvController, error) {
	if PvControllerInstance == nil {
		PvControllerInstance = &PvControllerImpl{}
	}
	return PvControllerInstance, nil
}

func (pc *PvControllerImpl) Run() {
	// 定期执行
	executor.ExecuteInPeriod(PvControllerDelay, PvControllerTimeGap, pc.syncPv)
}

func (pc *PvControllerImpl) syncPv() {
}

// AddPv 创建PersistentVolume
func (pc *PvControllerImpl) AddPv(pv *apiObject.PersistentVolume) error {
	// 将本地目录 /pvclient 挂载到服务器目录 /pvserver
	mountCmd := "mount " + config.NFSServer + ":/pvserver /pvclient"
	cmd := exec.Command("sh", "-c", mountCmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		return err
	}
	log.DebugLog("Bind to NFS server: " + string(output))

	// 在目录 /pvclient 创建目录 /:namespace/:name 作为PersistentVolume
	mkdirCmd := "mkdir -p /pvclient/" + pv.Metadata.Namespace + "/" + pv.Metadata.Name
	cmd = exec.Command("sh", "-c", mkdirCmd)
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		return err
	}
	log.DebugLog("Create PersistentVolume: " + string(output))

	return nil
}

// AddPvc 创建PersistentVolumeClaim
func (pc *PvControllerImpl) AddPvc(pvc *apiObject.PersistentVolumeClaim) error {
	return nil
}

// BindPvc 绑定PersistentVolumeClaim
func (pc *PvControllerImpl) BindPvc(pvc *apiObject.PersistentVolumeClaim) error {
	return nil
}
