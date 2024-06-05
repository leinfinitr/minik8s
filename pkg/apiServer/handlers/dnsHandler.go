package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiObject"
	etcdclient "minik8s/pkg/apiServer/etcdClient"
	"minik8s/pkg/config"
	httprequest "minik8s/tools/httpRequest"
	"minik8s/tools/log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetDNS(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if name == "" {
		log.ErrorLog("GetDNS name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	if namespace == "" {
		namespace = "default"
	}
	log.InfoLog("GetDNS: " + namespace + "/" + name)

	key := config.EtcdDnsPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("GetDNS: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(res) == 0 {
		log.ErrorLog("GetDNS: not found")
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	var resJson apiObject.Dns
	err = json.Unmarshal([]byte(res), &resJson)
	if err != nil {
		log.ErrorLog("GetDNS: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": resJson})
}

func GetDNSs(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	}
	log.InfoLog("GetDNSs: " + namespace)
	key := config.EtcdDnsPrefix + "/" + namespace
	res, err := etcdclient.EtcdStore.PrefixGet(key)
	if err != nil {
		log.ErrorLog("GetDNSs: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var dnss []apiObject.Dns
	for _, v := range res {
		var dns apiObject.Dns
		err = json.Unmarshal([]byte(v), &dns)
		if err != nil {
			log.ErrorLog("GetDNSs: " + err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		dnss = append(dnss, dns)
	}
	c.JSON(200, gin.H{"data": dnss})
}

func AddDNS(c *gin.Context) {
	// 在真正的服务之前，要确保是否已经创建出了Nginx的Pod
	nginxIP, err := GetNginxPod()
	if err != nil {
		log.ErrorLog("AddDNS: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 读取请求中的DNS对象
	var dns apiObject.Dns
	err = c.BindJSON(&dns)
	if err != nil {
		log.ErrorLog("AddDNS: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	log.InfoLog("AddDNS: " + dns.Metadata.Namespace + "/" + dns.Metadata.Name)
	dns.NginxIP = nginxIP

	// 检查DNS对象是否已经存在
	key := config.EtcdDnsPrefix + "/" + dns.Metadata.Namespace + "/" + dns.Metadata.Name
	oldRes, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("AddDNS: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(oldRes) != 0 {
		log.ErrorLog("AddDNS: already exists")
		c.JSON(409, gin.H{"error": "already exists"})
		return
	}
	if dns.Metadata.Name == "" {
		log.ErrorLog("AddDNS: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	if dns.Metadata.Namespace == "" {
		dns.Metadata.Namespace = "default"
	}

	// 获取每个子路径对应的service的IP
	for it, path := range dns.Spec.Paths {
		log.InfoLog("AddDNS: " + dns.Metadata.Namespace + "/" + dns.Metadata.Name + " path " + fmt.Sprint(it))
		if path.SvcName == "" {
			log.ErrorLog("AddDNS: svcName is empty")
			c.JSON(400, gin.H{"error": "svcName is empty"})
			return
		}
		svcKey := config.EtcdServicePrefix + "/" + dns.Metadata.Namespace + "/" + path.SvcName
		svcRes, err := etcdclient.EtcdStore.Get(svcKey)
		if err != nil {
			log.ErrorLog("AddDNS: " + err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if len(svcRes) == 0 {
			log.ErrorLog("AddDNS: svc not found")
			c.JSON(404, gin.H{"error": "svc not found"})
			return
		}
		service := apiObject.Service{}
		err = json.Unmarshal([]byte(svcRes), &service)
		if err != nil {
			log.ErrorLog("AddDNS: " + err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		dns.Spec.Paths[it].SvcIp = service.Spec.ClusterIP
	}
	dns.Metadata.UUID = uuid.New().String()

	// 更新每个节点的hosts文件
	Nodes := GetALLNodes()
	for _, node := range Nodes {
		url := "http://" + node.Status.Addresses[0].Address + ":" + fmt.Sprint(config.KubeproxyAPIPort) + config.DNSURI
		res, err := httprequest.PutObjMsg(url, dns)
		if err != nil || res.StatusCode != config.HttpSuccessCode {
			log.ErrorLog("AddDNS: " + err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	// 更新etcd
	dnsJson, err := json.Marshal(dns)
	if err != nil {
		log.ErrorLog("AddDNS: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Put(key, string(dnsJson))
	if err != nil {
		log.ErrorLog("AddDNS: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 在Nginx中增加相关配置
	if err = updateNginxConfig(&dns); err != nil {
		log.ErrorLog("AddDNS: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 更新覆盖Nginx的配置文件，并且重启Nginx
	if err = reloadNginx(); err != nil {
		log.ErrorLog("AddDNS: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": dns})

}

func DeleteDNS(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if name == "" {
		log.ErrorLog("DeleteDNS name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	if namespace == "" {
		namespace = "default"
	}
	log.InfoLog("DeleteDNS: " + namespace + "/" + name)

	key := config.EtcdDnsPrefix + "/" + namespace + "/" + name
	oldRes, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("DeleteDNS: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(oldRes) == 0 {
		log.ErrorLog("DeleteDNS: not found")
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	// 获取DNS对象
	var dns apiObject.Dns
	err = json.Unmarshal([]byte(oldRes), &dns)
	if err != nil {
		log.ErrorLog("DeleteDNS: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 删除每个节点的hosts文件
	Nodes := GetALLNodes()
	for _, node := range Nodes {
		url := "http://" + node.Status.Addresses[0].Address + ":" + fmt.Sprint(config.KubeproxyAPIPort) + config.DNSURI
		res, err := httprequest.DelMsg(url, dns)
		if err != nil || res.StatusCode != config.HttpSuccessCode {
			log.ErrorLog("DeleteDNS: " + err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	// 更新Nginx的配置文件
	if err = deleteNginxConfig(&dns); err != nil {
		log.ErrorLog("DeleteDNS: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 更新覆盖Nginx的配置文件，并且重启Nginx
	if err = reloadNginx(); err != nil {
		log.ErrorLog("DeleteDNS: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 删除etcd中的DNS对象
	err = etcdclient.EtcdStore.Delete(key)
	if err != nil {
		log.ErrorLog("DeleteDNS: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
}

func GetNginxPod() (string, error) {
	// 从etcd中获取Nginx的Pod，如果不存在则创建
	for {
		var nginxPod apiObject.Nginx
		res, err := etcdclient.EtcdStore.Get(config.EtcdNginxPrefix)
		if err == nil || res != "" {
			// Nginx的Pod已经存在
			err = json.Unmarshal([]byte(res), &nginxPod)
			if err == nil {
				if nginxPod.PodIP != "" {
					return nginxPod.PodIP, nil
				}
			}
		}
		if nginxPod.Phase == apiObject.PodRunning {
			return nginxPod.PodIP, nil
		}
		if nginxPod.Phase == "" {
			// 创建PVC
			//  go run minik8s/pkg/kubectl/main apply ./examples/dns/dns_pvc.yaml
			log.DebugLog("GetNginxPod: create pvc and create nginx pod")
			pvc := exec.Command("go", "run", "minik8s/pkg/kubectl/main", "apply", "./examples/dns/dns_pvc.yaml")
			err = pvc.Run()
			if err != nil {
				log.ErrorLog("Set PVC: " + err.Error())
				return "", err
			}
			time.Sleep(4 * time.Second)
			pod := exec.Command("go", "run", "minik8s/pkg/kubectl/main", "apply", "./examples/dns/dns_nginx.yaml")
			err = pod.Run()
			if err != nil {
				log.ErrorLog("Set Nginx: " + err.Error())
				return "", err
			}
			nginxPod.Phase = apiObject.PodBuilding
			resjson, err := json.Marshal(nginxPod)
			if err != nil {
				log.ErrorLog("GetNginxPod: " + err.Error())
				return "", err
			}
			err = etcdclient.EtcdStore.Put(config.EtcdNginxPrefix, string(resjson))
			if err != nil {
				log.ErrorLog("GetNginxPod: " + err.Error())
				return "", err
			}
		}

		time.Sleep(1 * time.Second)

	}
}

func updateNginxConfig(dns *apiObject.Dns) error {

	configPath := config.LocalConfigPath
	configPath = strings.Replace(configPath, ":namespace", dns.Metadata.Namespace, -1)
	configPath = strings.Replace(configPath, ":name", dns.Metadata.Name, -1)

	// 增加一个初始化的config到路径下
	err := exec.Command("cp", config.LocalOriginalPath, configPath).Run()
	if err != nil {
		log.ErrorLog("updateNginxConfig: " + err.Error())
		return err
	}

	// 读取文件
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 一行一行读取，在http块内的最后加上新的server
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// 先把最后的}删除，然后加上server之后再加上}
	if len(lines) > 0 {
		lines = lines[:len(lines)-1]
	}
	lines = append(lines, "    server {")
	lines = append(lines, "        listen 80;")
	lines = append(lines, "        server_name "+dns.Spec.Host+";")
	for _, path := range dns.Spec.Paths {
		lines = append(lines, "        location "+path.SubPath+" {")
		lines = append(lines, "            proxy_pass http://"+path.SvcIp+":"+path.SvcPort+"/;")
		lines = append(lines, "        }")
	}
	lines = append(lines, "    }")

	lines = append(lines, "}")

	// 写入文件
	err = os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		log.ErrorLog("updateNginxConfig: " + err.Error())
	}

	return nil
}

func reloadNginx() error {
	// 从etcd中获取Nginx的Pod
	var nginxPod apiObject.Nginx
	res, err := etcdclient.EtcdStore.Get(config.EtcdNginxPrefix)
	if err != nil || res == "" {
		return err
	}
	err = json.Unmarshal([]byte(res), &nginxPod)
	if err != nil {
		return err
	}

	url := "http://" + config.APIServerLocalAddress + ":" + fmt.Sprint(config.APIServerLocalPort) + config.PodExecURI
	url = strings.Replace(url, config.NameSpaceReplace, nginxPod.Namespace, -1)
	url = strings.Replace(url, config.NameReplace, nginxPod.Name, -1)
	url = strings.Replace(url, config.ContainerReplace, nginxPod.ContainerName, -1)

	cmd := apiObject.Command{Cmd: "cp /mnt/nginx.conf /etc/nginx/nginx.conf && nginx -s reload"}
	log.InfoLog("reloadNginx: " + url)
	resp, err := httprequest.PostObjMsg(url, cmd)
	if err != nil || resp.StatusCode != config.HttpSuccessCode {
		log.ErrorLog("reloadNginx: " + err.Error())
		return err
	}

	return nil
}

func deleteNginxConfig(dns *apiObject.Dns) error {
	configPath := config.LocalConfigPath
	configPath = strings.Replace(configPath, ":namespace", dns.Metadata.Namespace, -1)
	configPath = strings.Replace(configPath, ":name", dns.Metadata.Name, -1)
	// 读取文件
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 读取文件的每一行，遍历后找到需要被删除的server的前后下标，然后删除
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "server_name "+dns.Spec.Host) {
			// 把lines的最后两行删除
			if len(lines) > 1 {
				lines = lines[:len(lines)-2]
			}

			// 往后遍历寻找下一个server的位置
			for scanner.Scan() {
				line = scanner.Text()
				if strings.Contains(line, "server {") {
					lines = append(lines, line)
					break
				}
			}
		} else {
			lines = append(lines, line)
		}
	}

	// 写回文件
	err = os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		log.ErrorLog("deleteNginxConfig: " + err.Error())
	}

	return nil
}
