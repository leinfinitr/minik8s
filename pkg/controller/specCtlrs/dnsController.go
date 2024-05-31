package specctlrs

import (
	"encoding/json"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/executor"
	"minik8s/tools/log"
	"net/http"
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
				err = CreateDnsHandler(dns)
			case "Delete":
				err = DeleteDnsHandler(dns)
		}
		if err != nil {
			log.ErrorLog("syncDns: " + err.Error())
			return
		}
	}
	// 2. 获取所有的Service
	// 3. 更新DNS
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


func CreateDnsHandler(dns apiObject.Dns) error {
	return nil;
}

func DeleteDnsHandler(dns apiObject.Dns) error {
	// 1. 删除Dns
	// 2. 创建DnsRequest
	return nil
}