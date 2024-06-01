package specctlrs

import (
	"encoding/json"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/executor"
	"minik8s/tools/host"
	httprequest "minik8s/tools/httpRequest"
	"minik8s/tools/log"
	netRequest "minik8s/tools/netRequest"
	"minik8s/tools/nginx"
	stringops "minik8s/tools/stringops"
	"net/http"
	"os"
	"strings"
	"time"
)

type DnsController interface {
	Run()
}

type DnsControllerImpl struct {
	hostList     []string
	nginxSvcName string
	nginxSvcIp   string
}

func NewDnsController() (DnsController, error) {
	return &DnsControllerImpl{}, nil
}

var(
	DnsControllerDelay   = 0 * time.Second
	DnsControllerTimeGap = []time.Duration{5 * time.Second}
)

func (dc *DnsControllerImpl) Run() {
	dc.CreateNginx()
	dc.UpdateNginxIp()
	// 定期执行
	executor.ExecuteInPeriod(DnsControllerDelay, DnsControllerTimeGap, dc.syncDns)
}

func (dc *DnsControllerImpl) syncDns() {
	// 1. 获取所有的DnsRequest
	dnsRequests, err := GetAllDnsRequest()
	if err != nil {
		log.ErrorLog("syncDns: " + err.Error())
		return
	}
	for _, dnsRequest := range dnsRequests {
		var dns apiObject.Dns
		namespace := dnsRequest.DnsMeta.Namespace
		name := dnsRequest.DnsMeta.Name
		url := config.APIServerURL() + config.DnsRequestURI
		url = strings.Replace(url, config.NameSpaceReplace, namespace, -1)
		url = strings.Replace(url, config.NameReplace, name, -1)
		res, err := http.Get(url)
		if err != nil {
			log.ErrorLog("syncDns: " + err.Error())
			return
		}
		err = json.NewDecoder(res.Body).Decode(&dns)
		requestType := dnsRequest.Action
		switch requestType {
			case "Create":
				err = dc.CreateDnsHandler(dns)
			case "Delete":
				err = dc.DeleteDnsHandler(dns)
			case "Update":
				err = dc.DeleteDnsHandler(dns)
				if err != nil {
					log.ErrorLog("syncDns: " + err.Error())
					return
				}
				err = dc.CreateDnsHandler(dns)
		}
		if err != nil {
			log.ErrorLog("syncDns: " + err.Error())
			return
		}
		_,err = httprequest.DelMsg(url,dns)
		if err != nil {
			log.ErrorLog("syncDns: " + err.Error())
			return
		}
	}
}

func GetAllDnsRequest() (dnsRequests []apiObject.DnsRequest, err error) {
	// 1. 获取所有的DnsRequest
	url := config.APIServerURL() +config.GlobalDnsRequestURI
	res,err := http.Get(url)
	if err != nil{
		log.ErrorLog("GetAllDnsRequest: "+err.Error())
		return dnsRequests,err
	}
	if res.StatusCode != 200{
		log.ErrorLog("GetAllDnsRequest: "+res.Status)
		return dnsRequests,err
	}
	err = json.NewDecoder(res.Body).Decode(&dnsRequests)
	if err != nil{
		log.ErrorLog("GetAllDnsRequest: "+err.Error())
		return dnsRequests,err
	}
	return dnsRequests,nil
}


func (dc *DnsControllerImpl)CreateDnsHandler(dns apiObject.Dns) error {
	if dns.Spec.Host == "" {
		log.ErrorLog("CreateDnsHandler: Host is empty")
		return nil
	}
	if dns.Metadata.Namespace == "" {
		dns.Metadata.Namespace = "default"
	}
	nginxConf := nginx.TranslateConfig(dns)
	hostEntry := dc.nginxSvcIp + " " + dns.Spec.Host
	dc.hostList = append(dc.hostList, hostEntry)
	hostRequest := apiObject.HostRequest{
		Action:  "Create",
		DnsObject: dns,
		DnsConfig: nginxConf,
		HostList: dc.hostList,
	}
	err := BroadcastProxy(hostRequest)
	if err != nil {
		log.ErrorLog("CreateDnsHandler: " + err.Error())
		return err
	}
	return nil
}

func (dc *DnsControllerImpl)DeleteDnsHandler(dns apiObject.Dns) error {
	// 1. 删除Dns
	// 2. 创建DnsRequest
	return nil
}

func BroadcastProxy(hostRequest apiObject.HostRequest) error {
	//TODO
	// 1. 获取所有的Proxy
	// 2. 广播HostRequest
	return nil
}

func (dc *DnsControllerImpl)CreateNginx() {
	filePath := config.NginxPodYamlPath
	content,err := os.ReadFile(filePath)
	if err != nil {
		log.ErrorLog("CreateNginx: "+err.Error())
		return
	}
	nginxPod := apiObject.Pod{}
	err = json.Unmarshal(content,&nginxPod)
	if err != nil {
		log.ErrorLog("CreateNginx: "+err.Error())
		return
	}
	if nginxPod.Metadata.Namespace == "" {
		nginxPod.Metadata.Namespace = "default"
	}
	url := config.APIServerURL() + config.PodsURI
	url = strings.Replace(url,config.NameSpaceReplace,nginxPod.Metadata.Namespace,-1)
	nginxPod.Metadata.Name += "-"+stringops.GenerateRandomString(5)
	nginxPod.Spec.NodeName,err = host.GetHostname()
	if err != nil {
		log.ErrorLog("CreateNginx: "+err.Error())
		return
	}
	code,_,err := netRequest.PostRequestByTarget(url,&nginxPod)
	if err != nil {
		log.ErrorLog("CreateNginx: "+err.Error())
		return
	}
	if code != http.StatusCreated {
		log.ErrorLog("CreateNginx: code is not 201")
		return
	}
	log.InfoLog("CreateNginx: success")
}

func (dc *DnsControllerImpl)UpdateNginxIp() {
	//TODO
	// 1. 获取Nginx的SvcName和SvcIp
	// 2. 更新本地Nginx的SvcName和SvcIp
}