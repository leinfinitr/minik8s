package specctlrs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/executor"
	"minik8s/tools/log"

	etcdclient "minik8s/pkg/apiServer/etcdClient"
)

type PvController interface {
	Run()
}

type PvControllerImpl struct {
	// 服务器地址
	Address string
	// 服务器端口
	Port int
	// 转发请求
	Router *gin.Engine
	// 用于存储PersistentVolume，名称为namespace/name
	PvMap map[string]*apiObject.PersistentVolume
	// 用于存储PersistentVolumeClaim和PersistentVolume的映射关系，名称均为为namespace/name
	PvcPvMap map[string]string
}

var (
	PvControllerDelay   = 3 * time.Second
	PvControllerTimeGap = []time.Duration{10 * time.Second}
)

func NewPvController() (PvController, error) {
	// 设置gin的运行模式
	gin.SetMode(gin.ReleaseMode)

	return &PvControllerImpl{
		Address:  config.PVServerAddress,
		Port:     config.PVServerPort,
		Router:   gin.New(),
		PvMap:    make(map[string]*apiObject.PersistentVolume),
		PvcPvMap: make(map[string]string),
	}, nil
}

func (pc *PvControllerImpl) Run() {
	// 注册路由
	pc.Register()

	// 开启一个协程定期执行同步函数
	go executor.ExecuteInPeriod(PvControllerDelay, PvControllerTimeGap, pc.syncPv)

	// 开启线程用于处理请求
	err := pc.Router.Run(pc.Address + ":" + fmt.Sprint(pc.Port))
	if err != nil {
		log.ErrorLog("ServerlessServer Run: " + err.Error())
	}

}

// Register 注册路由
func (pc *PvControllerImpl) Register() {
	// 创建PersistentVolume
	pc.Router.POST(config.PersistentVolumesURI, pc.CreatePv)

	// 创建PersistentVolumeClaim
	pc.Router.POST(config.PersistentVolumeClaimsURI, pc.CreatePvc)

	// 绑定pod到PersistentVolumeClaim
	pc.Router.POST(config.PersistentVolumeClaimURI, pc.BindPodToPvc)
	// 解绑pod和PersistentVolumeClaim
	pc.Router.DELETE(config.PersistentVolumeClaimURI, pc.UnbindPodToPvc)
	// 获取PersistentVolumeClaim绑定的PersistentVolume
	pc.Router.GET(config.PersistentVolumeClaimURI, pc.GetPvcBind)
}

// CreatePv 创建PersistentVolume
func (pc *PvControllerImpl) CreatePv(c *gin.Context) {
	pv := &apiObject.PersistentVolume{}
	err := c.ShouldBindJSON(pv)
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	err = pc.addPv(pv)
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": "Create PersistentVolume " + pv.Metadata.Name})
}

// CreatePvc 创建PersistentVolumeClaim
func (pc *PvControllerImpl) CreatePvc(c *gin.Context) {
	pvc := &apiObject.PersistentVolumeClaim{}
	err := c.ShouldBindJSON(pvc)
	if err != nil {
		log.ErrorLog("Create PersistentVolumeClaim: " + err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	err = pc.addPvc(pvc)
	if err != nil {
		log.ErrorLog("Create PersistentVolumeClaim: " + err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": "Create PersistentVolumeClaim " + pvc.Metadata.Name})
}

// BindPodToPvc 绑定PersistentVolumeClaim
func (pc *PvControllerImpl) BindPodToPvc(c *gin.Context) {
	podNamespace := c.Param("namespace")
	podName := c.Param("name")
	log.DebugLog("Bind PersistentVolumeClaim: " + podNamespace + "/" + podName)

	pvc := &apiObject.PersistentVolumeClaim{}
	err := c.ShouldBindJSON(pvc)
	if err != nil {
		log.ErrorLog("Bind PersistentVolumeClaim: " + err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	err = pc.bindPodToPvc(pvc, podNamespace+"/"+podName)
	if err != nil {
		log.ErrorLog("Bind PersistentVolumeClaim: " + err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Bind PersistentVolumeClaim " + pvc.Metadata.Name})
}

// UnbindPodToPvc 解绑PersistentVolumeClaim
func (pc *PvControllerImpl) UnbindPodToPvc(c *gin.Context) {
	pvcNamespace := c.Param("namespace")
	pvcName := c.Param("name")
	key := config.EtcdPvcPrefix + "/" + pvcNamespace + "/" + pvcName
	pvc := &apiObject.PersistentVolumeClaim{}
	response, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("Unbind PersistentVolumeClaim: " + err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	err = json.Unmarshal([]byte(response), pvc)
	if err != nil {
		log.ErrorLog("Unbind PersistentVolumeClaim: " + err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	err = pc.unbindPodToPvc(pvc)
	if err != nil {
		log.ErrorLog("Unbind PersistentVolumeClaim: " + err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Unbind PersistentVolumeClaim " + pvc.Metadata.Name})
}

// GetPvcBind 获取PersistentVolumeClaim绑定的PersistentVolume
func (pc *PvControllerImpl) GetPvcBind(c *gin.Context) {
	pvcNamespace := c.Param("namespace")
	pvcName := c.Param("name")
	pvName := pc.getPvcBind(pvcNamespace + "/" + pvcName)
	c.JSON(http.StatusOK, pvName)
}

// syncPv 同步PersistentVolume
func (pc *PvControllerImpl) syncPv() {
	// 从etcd中获取所有PersistentVolumeClaim
	response, err := etcdclient.EtcdStore.PrefixGet(config.EtcdPvcPrefix)
	if err != nil {
		log.ErrorLog("Sync PersistentVolume: " + err.Error())
		return
	}
	// 绑定PersistentVolumeClaim
	for _, v := range response {
		pvc := apiObject.PersistentVolumeClaim{}
		err = json.Unmarshal([]byte(v), &pvc)
		if err != nil {
			log.ErrorLog("Sync PersistentVolume: " + err.Error())
			continue
		}
		if pvc.Status.Phase == apiObject.ClaimPending {
			err = pc.bindPvcToPv(&pvc)
			if err != nil {
				log.ErrorLog("Sync PersistentVolume: " + err.Error())
				continue
			}
		}
	}
}

// addPv 创建PersistentVolume
func (pc *PvControllerImpl) addPv(pv *apiObject.PersistentVolume) error {
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
	err = cmd.Run()
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		return err
	}
	log.DebugLog("Bind to NFS server: " + config.NFSServer + ":" + config.PVServerPath)
	// 在目录 /pvclient 创建目录 /:namespace/:name 作为PersistentVolume
	mkdirCmd := "mkdir -p " + config.PVClientPath + "/" + pv.Metadata.Namespace + "/" + pv.Metadata.Name
	cmd = exec.Command("sh", "-c", mkdirCmd)
	err = cmd.Run()
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		return err
	}
	log.DebugLog("Create PersistentVolume: " + pvNamespace + "/" + pvName)
	// 清空目录 /pvclient/:namespace/:name
	rmCmd := "rm -rf " + config.PVClientPath + "/" + pv.Metadata.Namespace + "/" + pv.Metadata.Name + "/*"
	cmd = exec.Command("sh", "-c", rmCmd)
	err = cmd.Run()
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		return err
	}
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

// addPvc 创建PersistentVolumeClaim
func (pc *PvControllerImpl) addPvc(pvc *apiObject.PersistentVolumeClaim) error {
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
	err = pc.bindPvcToPv(pvc)
	if err != nil {
		log.ErrorLog("Create PersistentVolumeClaim: " + err.Error())
		return err
	}

	return nil
}

// getPvcBind 获取PersistentVolumeClaim绑定的PersistentVolume
func (pc *PvControllerImpl) getPvcBind(pvcName string) string {
	return pc.PvcPvMap[pvcName]
}

// bindPodToPvc 绑定Pod到PersistentVolumeClaim
func (pc *PvControllerImpl) bindPodToPvc(pvc *apiObject.PersistentVolumeClaim, podName string) error {
	// 修改pvc的状态
	pvc.Status.IsBound = true
	pvc.Status.BoundPodName = podName
	// 更新pvc的状态
	err := pc.updatePvc(pvc)
	if err != nil {
		log.ErrorLog("Bind Pod to PersistentVolumeClaim: " + err.Error())
		return err
	}

	log.InfoLog("Bind Pod to PersistentVolumeClaim: " + pvc.Metadata.Namespace + "/" + pvc.Metadata.Name + " bound to " + podName)
	return nil
}

// unbindPodToPvc 解绑Pod和PersistentVolumeClaim
func (pc *PvControllerImpl) unbindPodToPvc(pvc *apiObject.PersistentVolumeClaim) error {
	// 修改pvc的状态
	pvc.Status.IsBound = false
	pvc.Status.BoundPodName = ""
	// 更新pvc的状态
	err := pc.updatePvc(pvc)
	if err != nil {
		log.ErrorLog("Unbind Pod to PersistentVolumeClaim: " + err.Error())
		return err
	}

	log.InfoLog("Unbind Pod to PersistentVolumeClaim: " + pvc.Metadata.Namespace + "/" + pvc.Metadata.Name)
	return nil
}

// bindPvcToPv 绑定PersistentVolumeClaim
func (pc *PvControllerImpl) bindPvcToPv(pvc *apiObject.PersistentVolumeClaim) error {
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
	err = pc.updatePvc(pvc)
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

// updatePvc 在etcd中更新PersistentVolumeClaim
func (pc *PvControllerImpl) updatePvc(pvc *apiObject.PersistentVolumeClaim) error {
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

// newPV 通过Pvc请求自动创建一个PersistentVolume
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
	err := pc.addPv(pv)
	if err != nil {
		log.ErrorLog("New PersistentVolume: " + err.Error())
		return nil
	}
	return pv
}
