package specctlrs

import (
	"encoding/json"
	"os/exec"
	"time"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/executor"
	"minik8s/tools/log"

	etcdclient "minik8s/pkg/apiServer/etcdClient"
)

type PvController interface {
	Run()
	AddPv(pv *apiObject.PersistentVolume) error
	AddPvc(pvc *apiObject.PersistentVolumeClaim) error
	GetBind(pvcName string) string
}

type PvControllerImpl struct {
	// 用于存储PersistentVolume
	PvMap map[string]*apiObject.PersistentVolume
	// 用于存储PersistentVolumeClaim和PersistentVolume的映射关系
	PvcPvMap map[string]string
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

// syncPv 同步PersistentVolume
func (pc *PvControllerImpl) syncPv() {
	for {
		// 从etcd中获取所有PersistentVolumeClaim
		response, err := etcdclient.EtcdStore.Get(config.EtcdPvcPrefix)
		if err != nil {
			log.ErrorLog("Sync PersistentVolume: " + err.Error())
			continue
		}
		// 绑定PersistentVolumeClaim
		var pvcList []apiObject.PersistentVolumeClaim
		err = json.Unmarshal([]byte(response), &pvcList)
		if err != nil {
			log.ErrorLog("Sync PersistentVolume: " + err.Error())
			continue
		}
		for _, pvc := range pvcList {
			if pvc.Status.Phase == apiObject.ClaimPending {
				err = pc.bindPvc(&pvc)
				if err != nil {
					log.ErrorLog("Sync PersistentVolume: " + err.Error())
					continue
				}
			}
		}
	}
}

// AddPv 创建PersistentVolume
func (pc *PvControllerImpl) AddPv(pv *apiObject.PersistentVolume) error {
	// 检查pv是否已经存在
	pvName := pv.Metadata.Name
	pvNamespace := pv.Metadata.Namespace
	key := config.EtcdPvPrefix + "/" + pvNamespace + "/" + pvName
	response, err := etcdclient.EtcdStore.Get(key)
	if response != "" {
		log.ErrorLog("Create PersistentVolume: pv already exists" + response)
		return err
	}
	// 将本地目录 /pvclient 挂载到服务器目录 /pvserver
	mountCmd := "mount " + config.NFSServer + ":" + config.PVServerPath + " " + config.PVClientPath
	cmd := exec.Command("sh", "-c", mountCmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		return err
	}
	log.DebugLog("Bind to NFS server: " + string(output))
	// 在目录 /pvclient 创建目录 /:namespace/:name 作为PersistentVolume
	mkdirCmd := "mkdir -p " + config.PVClientPath + "/" + pv.Metadata.Namespace + "/" + pv.Metadata.Name
	cmd = exec.Command("sh", "-c", mkdirCmd)
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		return err
	}
	log.DebugLog("Create PersistentVolume: " + string(output))
	// 修改pv的状态
	pv.Status.Phase = apiObject.VolumeAvailable
	// 将pv存入map
	pc.PvMap[pv.Metadata.Namespace+"/"+pv.Metadata.Name] = pv
	// 将pv存入etcd
	pvJson, err := json.Marshal(&pv)
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		return err
	}
	err = etcdclient.EtcdStore.Put(key, string(pvJson))
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		return err
	}

	return nil
}

// AddPvc 创建PersistentVolumeClaim
func (pc *PvControllerImpl) AddPvc(pvc *apiObject.PersistentVolumeClaim) error {
	// 检查pvc是否已经存在
	pvcName := pvc.Metadata.Name
	pvcNamespace := pvc.Metadata.Namespace
	key := config.EtcdPvcPrefix + "/" + pvcNamespace + "/" + pvcName
	response, err := etcdclient.EtcdStore.Get(key)
	if response != "" {
		log.ErrorLog("Create PersistentVolumeClaim: pvc already exists" + response)
		return err
	}
	log.DebugLog("Create PersistentVolumeClaim: " + pvcNamespace + "/" + pvcName)
	// 修改pvc的状态
	pvc.Status.Phase = apiObject.ClaimPending
	// 将pvc存入map，其所对应的pv为nil
	pc.PvcPvMap[pvc.Metadata.Namespace+"/"+pvc.Metadata.Name] = ""
	// 将pvc存入etcd
	pvcJson, err := json.Marshal(pvc)
	if err != nil {
		log.ErrorLog("Create PersistentVolumeClaim: " + err.Error())
		return err

	}
	err = etcdclient.EtcdStore.Put(key, string(pvcJson))
	if err != nil {
		log.ErrorLog("Create PersistentVolumeClaim: " + err.Error())
		return err
	}
	// 主动绑定pvc
	err = pc.bindPvc(pvc)
	if err != nil {
		log.ErrorLog("Create PersistentVolumeClaim: " + err.Error())
		return err
	}

	return nil
}

// GetBind 获取PersistentVolumeClaim绑定的PersistentVolume
func (pc *PvControllerImpl) GetBind(pvcName string) string {
	return pc.PvcPvMap[pvcName]
}

// bindPvc 绑定PersistentVolumeClaim
func (pc *PvControllerImpl) bindPvc(pvc *apiObject.PersistentVolumeClaim) error {
	pvcKey := pvc.Metadata.Namespace + "/" + pvc.Metadata.Name
	var pv *apiObject.PersistentVolume
	// 从PvMap中找到一个未绑定的pv，并绑定
	for _, v := range pc.PvMap {
		if v.Status.Phase == apiObject.VolumeAvailable {
			pv = v
			break
		}
	}
	// 若没有找到合适的pv，则创建一个新的pv
	if pv == nil {
		pv = pc.newPV(pvc)
	}
	// 若创建pv失败，则返回错误
	if pv == nil {
		log.ErrorLog("Bind PersistentVolumeClaim: create new pv failed")
		return nil
	}
	// 将pvc绑定到pv
	pvKey := pv.Metadata.Namespace + "/" + pv.Metadata.Name
	pc.PvcPvMap[pvcKey] = pvKey
	log.DebugLog("Bind PersistentVolumeClaim: " + pvcKey + " to " + pvKey)
	// 更新pv的状态
	pv.Status.Phase = apiObject.VolumeBound
	err := pc.updatePv(pv)
	if err != nil {
		log.ErrorLog("Bind PersistentVolumeClaim: " + err.Error())
		return err
	}
	// 更新pvc的状态
	pvc.Status.Phase = apiObject.ClaimBound
	err = pc.UpdatePvc(pvc)
	if err != nil {
		log.ErrorLog("Bind PersistentVolumeClaim: " + err.Error())
		return err
	}

	log.InfoLog("Bind PersistentVolumeClaim: " + pvcKey + " bound to " + pvKey)
	return nil
}

// updatePv 在etcd中更新PersistentVolume
func (pc *PvControllerImpl) updatePv(pv *apiObject.PersistentVolume) error {
	pvName := pv.Metadata.Name
	pvNamespace := pv.Metadata.Namespace
	key := config.EtcdPvPrefix + "/" + pvNamespace + "/" + pvName
	pvJson, err := json.Marshal(pv)
	if err != nil {
		log.ErrorLog("Update PersistentVolume status: " + err.Error())
		return err
	}
	err = etcdclient.EtcdStore.Put(key, string(pvJson))
	if err != nil {
		log.ErrorLog("Update PersistentVolume status: " + err.Error())
		return err
	}
	return nil
}

// UpdatePvc 在etcd中更新PersistentVolumeClaim
func (pc *PvControllerImpl) UpdatePvc(pvc *apiObject.PersistentVolumeClaim) error {
	pvcName := pvc.Metadata.Name
	pvcNamespace := pvc.Metadata.Namespace
	key := config.EtcdPvcPrefix + "/" + pvcNamespace + "/" + pvcName
	pvcJson, err := json.Marshal(pvc)
	if err != nil {
		log.ErrorLog("Update PersistentVolumeClaim status: " + err.Error())
		return err
	}
	err = etcdclient.EtcdStore.Put(key, string(pvcJson))
	if err != nil {
		log.ErrorLog("Update PersistentVolumeClaim status: " + err.Error())
		return err
	}
	return nil
}

// newPV 通过Pvc请求创建一个PersistentVolume
func (pc *PvControllerImpl) newPV(pvc *apiObject.PersistentVolumeClaim) *apiObject.PersistentVolume {
	// 创建一个PersistentVolume
	pv := &apiObject.PersistentVolume{
		TypeMeta: apiObject.TypeMeta{
			Kind:       "PersistentVolume",
			APIVersion: "v1",
		},
		Metadata: apiObject.ObjectMeta{
			Name:      pvc.Metadata.Name,
			Namespace: pvc.Metadata.Namespace,
		},
		Spec: apiObject.PersistentVolumeSpec{
			Capacity:      pvc.Spec.Resources,
			AccessModes:   pvc.Spec.AccessModes,
			ReclaimPolicy: apiObject.Recycle,
			Remote: apiObject.NetworkFileSystem{
				Server: config.NFSServer,
				Path:   "/",
			},
		},
		Status: apiObject.PersistentVolumeStatus{
			Phase: apiObject.VolumePending,
		},
	}
	// 创建PersistentVolume
	err := pc.AddPv(pv)
	if err != nil {
		log.ErrorLog("New PersistentVolume: " + err.Error())
		return nil
	}
	return pv
}
