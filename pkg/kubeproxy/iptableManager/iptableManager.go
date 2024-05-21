package iptableManager

import (
	"fmt"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/entity"
	"minik8s/tools/log"

	"github.com/coreos/go-iptables/iptables"
)

type IptableManager interface {
	CreateService(createEvent *entity.ServiceEvent) error
	UpdateService(updateEvent *entity.ServiceEvent) error
	DeleteService(deleteEvent *entity.ServiceEvent) error
}

type iptableManager struct {
	// 实现从service的name到Endpoint信息的映射
	service2Endpoint map[string][]apiObject.Endpoint

	// 用来存储iptables的规则
	iptable *iptables.IPTables
}

var iptableMgr *iptableManager

func GetIptableManager() IptableManager {
	if iptableMgr == nil {
		iptableMgr = &iptableManager{
			service2Endpoint: make(map[string][]apiObject.Endpoint),
		}
	}

	iptableMgr.init_iptable()

	return iptableMgr
}

// 在此处，proxy主要是修改iptables的NAT表，实现service的有关功能
func (i *iptableManager) init_iptable() {
	log.InfoLog("[KUBEPROXY]: init iptables")

	i.iptable, _ = iptables.New()

	//设置NAT表的策略
	i.iptable.ChangePolicy("nat", "PREROUTING", "ACCEPT")
	i.iptable.ChangePolicy("nat", "INPUT", "ACCEPT")
	i.iptable.ChangePolicy("nat", "OUTPUT", "ACCEPT")
	i.iptable.ChangePolicy("nat", "POSTROUTING", "ACCEPT")

	// 创建 NAT 表中的新链
	i.iptable.NewChain("nat", "KUBE-SERVICES")
	i.iptable.NewChain("nat", "KUBE-POSTROUTING")
	i.iptable.NewChain("nat", "KUBE-MARK-MASQ")
	i.iptable.NewChain("nat", "KUBE-NODEPORTS")

	/* 往 NAT 表新创建的链中添加规则 */
	// 从 PREROUTING 链中转发到 KUBE-SERVICES 链，负责非本机访问
	i.iptable.Append("nat", "PREROUTING", "-j", "KUBE-SERVICES", "-m", "comment", "--comment", "kubernetes service portals")
	// 从 OUTPUT 链中转发到 KUBE-SERVICES 链，负责本机访问
	i.iptable.Append("nat", "OUTPUT", "-j", "KUBE-SERVICES", "-m", "comment", "--comment", "kubernetes service portals")

	i.iptable.Append("nat", "POSTROUTING", "-j", "KUBE-POSTROUTING", "-m", "comment", "--comment", "kubernetes postrouting rules")
	i.iptable.AppendUnique("nat", "KUBE-MARK-MASQ", "-j", "MARK", "--or-mark", "0x4000")
	i.iptable.AppendUnique("nat", "KUBE-POSTROUTING", "-m", "comment", "--comment", "kubernetes service traffic requiring SNAT", "-j", "MASQUERADE", "-m", "mark", "--mark", "0x4000/0x4000")
}

// 本地访问时的规则
// PREROUTING --> KUBE-SERVICE --> KUBE-NODEPORTS --> KUBE-SVC-XXX --> KUBE-SEP-XXX
// 外部访问时的规则
// OUTPUT --> KUBE-SERVICE --> KUBE-NODEPORTS --> KUBE-SVC-XXX --> KUBE-SEP-XXX

func (i *iptableManager) CreateService(createEvent *entity.ServiceEvent) error {
	log.InfoLog("[KUBEPROXY]: Start CreateService")
	serviceName := createEvent.Service.Metadata.Name
	var podsIPList []string
	for _, endpoint := range createEvent.Endpoints {
		podsIPList = append(podsIPList, endpoint.IP)
	}

	// 根据service制定开放的port端口，分别设置规则
	if createEvent.Service.Spec.Type == "ClusterIP" || createEvent.Service.Spec.Type == "" {
		clusterIP := createEvent.Service.Spec.ClusterIP
		for _, portConfig := range createEvent.Service.Spec.Ports {
			msg := fmt.Sprintf("[KUBEPROXY]: Start CreateService Port :%d", portConfig.Port)
			log.InfoLog(msg)
			port := portConfig.Port
			targetPort := portConfig.TargetPort
			protocol := portConfig.Protocol
			err := i.setIPTableForClusterIP(serviceName, clusterIP, port, targetPort, protocol, podsIPList)
			if err != nil {
				log.ErrorLog("CreateService error: " + err.Error())
				return err
			}
		}
		// TODO: 需要考虑不指定端口的情况

	} else if createEvent.Service.Spec.Type == "NodePort" {
		for _, portConfig := range createEvent.Service.Spec.Ports {
			port := portConfig.Port
			targetPort := portConfig.TargetPort
			nodePort := portConfig.NodePort
			protocol := portConfig.Protocol
			err := i.setIPtableForNodePort(serviceName, "", port, targetPort, nodePort, protocol, podsIPList)
			if err != nil {
				log.ErrorLog("CreateService error: " + err.Error())
				return err
			}
		}
	} else {
		log.ErrorLog("CreateService error: service type error")
		return fmt.Errorf("CreateService error: service type error")
	}

	// 把service的UUID和endpoint信息存储起来
	for _, endpoint := range createEvent.Endpoints {
		i.service2Endpoint[createEvent.Service.Metadata.Name] = append(i.service2Endpoint[createEvent.Service.Metadata.UUID], endpoint)
	}
	return nil
}

func (i *iptableManager) setIPTableForClusterIP(serviceName string, clusterIP string, port int32, targetPort int32, protocol apiObject.Protocol, podsIPList []string) error {
	log.InfoLog("[KUBEPROXY]: serviceName:" + serviceName + " clusterIP:" + clusterIP + " port:" + fmt.Sprintf("%d", port) + " targetPort:" + fmt.Sprintf("%d", targetPort))

	// 为该Service Port创建一个独有的KUBE—SVC链
	kubeSvc := "KUBE-SVC-" + serviceName + "-" + fmt.Sprintf("%d", port)
	err := i.iptable.NewChain("nat", kubeSvc)
	if err != nil {
		log.ErrorLog("KUBE-SVC NewChain error: " + err.Error())
		return err
	}

	// 在KUBE-SERVICES链中添加规则转发至KUBE-SVC链
	err = i.iptable.Insert("nat", "KUBE-SERVICES", 1, "-d", clusterIP, "-p", string(protocol), "--dport", fmt.Sprintf("%d", port), "-j", kubeSvc)
	if err != nil {
		log.ErrorLog("KUBE-SERVICES Append error: " + err.Error())
		return err
	}

	// 在KUBE-SERVICES链中添加规则，标记数据包
	err = i.iptable.Insert("nat", "KUBE-SERVICES", 1, "-d", clusterIP, "-p", string(protocol), "--dport", fmt.Sprintf("%d", port), "-j", "KUBE-MARK-MASQ")
	if err != nil {
		log.ErrorLog("KUBE-SERVICES Mark error: " + err.Error())
		return err
	}

	// 在KUBE-SVC链中添加规则，转发至对应的pod
	podsLen := len(podsIPList)
	for j, podIP := range podsIPList {
		// 为每个pod创建一个独有的KUBE-SEP链
		kubeSep := "KUBE-SEP-" + serviceName + "-" + fmt.Sprintf("%d", port) + "-" + fmt.Sprintf("%d", j)
		err = i.iptable.NewChain("nat", kubeSep)
		if err != nil {
			log.ErrorLog("KUBE-SEP NewChain error: " + err.Error())
			return err
		}

		// 在KUBE-SVC链中添加规则转发至KUBE-SEP链，默认采用随机算法
		if j == podsLen-1 {
			err = i.iptable.Append("nat", kubeSvc, "-j", kubeSep)
			if err != nil {
				log.ErrorLog("KUBE-SVC forward to KUBE_SEP error: " + err.Error())
				return err
			}
		} else {
			Probability := 1.0 / float64(podsLen-j)
			err = i.iptable.Append("nat", kubeSvc, "-m", "statistic", "--mode", "random", "--probability", fmt.Sprintf("%f", Probability), "-j", kubeSep)
			if err != nil {
				log.ErrorLog("KUBE-SVC forward to KUBE_SEP error: " + err.Error())
				return err
			}
		}

		// 在KUBE-SEP链中添加规则，转发至对应的pod
		err = i.iptable.Append("nat", kubeSep, "-d", clusterIP, "-p", string(protocol), "--dport", fmt.Sprintf("%d", port), "-j", "DNAT", "--to-destination", podIP+":"+fmt.Sprintf("%d", targetPort))
		if err != nil {
			log.ErrorLog("KUBE-SEP DNAT error: " + err.Error())
			return err
		}

		// 在KUBE-SEP链中添加规则，标记数据包
		err = i.iptable.Insert("nat", kubeSep, 1, "-d", clusterIP, "-p", string(protocol), "--dport", fmt.Sprintf("%d", port), "-j", "KUBE-MARK-MASQ")
		if err != nil {
			log.ErrorLog("KUBE-SEP Mark error: " + err.Error())
			return err
		}
	}

	return nil
}

// 注意，NodePort的实现和ClusterIP的实现存在差异，会多出一条链
func (i *iptableManager) setIPtableForNodePort(serviceName string, clusterIP string, port int32, targetPort int32, nodePort int32, protocol apiObject.Protocol, podsIPList []string) error {
	log.InfoLog("[KUBEPROXY]: NodePort  serviceName:" + serviceName + " clusterIP:" + clusterIP + " port:" + fmt.Sprintf("%d", port) + " targetPort:" + fmt.Sprintf("%d", targetPort) + " nodePort:" + fmt.Sprintf("%d", nodePort))

	// 为该Service Port创建一个独有的KUBE—SVC链
	kubeSvc := "KUBE-SVC-" + serviceName + "-" + fmt.Sprintf("%d", port)
	err := i.iptable.NewChain("nat", kubeSvc)
	if err != nil {
		log.ErrorLog("KUBE-SVC NewChain error: " + err.Error())
		return err
	}

	// 在KUBE-SERVICES链末尾添加规则，转发至KUBE-NODEPORTS链
	err = i.iptable.AppendUnique("nat", "KUBE-SERVICES", "-m", "Kubenetes service nodeport; NOTE: this must be the last rule in this chain",
		"-p", string(protocol), "-j", "KUBE-NODEPORTS")
	if err != nil {
		log.ErrorLog("KUBE-SERVICES Append error: " + err.Error())
		return err
	}

	// 在KUBE-NODEPORTS链中添加规则，转发至KUBE-SVC链
	err = i.iptable.Insert("nat", "KUBE-NODEPORTS", 1, "-p", string(protocol), "--dport", fmt.Sprintf("%d", nodePort), "-j", kubeSvc)
	if err != nil {
		log.ErrorLog("KUBE-NODEPORTS Append error: " + err.Error())
		return err
	}

	// 在KUBE-NODEPORTS链中添加规则，标记数据包
	err = i.iptable.Insert("nat", "KUBE-NODEPORTS", 1, "-p", string(protocol), "--dport", fmt.Sprintf("%d", nodePort), "-j", "KUBE-MARK-MASQ")
	if err != nil {
		log.ErrorLog("KUBE-NODEPORTS Mark error: " + err.Error())
		return err
	}

	// 在KUBE-SVC链的上一条规则的前面添加规则，标记数据包
	err = i.iptable.Insert("nat", kubeSvc, 1, "-d", clusterIP, "-p", string(protocol), "--dport", fmt.Sprintf("%d", port), "-j", "KUBE-MARK-MASQ")
	if err != nil {
		log.ErrorLog("KUBE-SVC Mark error: " + err.Error())
		return err
	}

	// 在KUBE-SVC链中添加规则，转发至对应的pod
	podsLen := len(podsIPList)
	for j, podIP := range podsIPList {
		// 为每个pod创建一个独有的KUBE-SEP链
		kubeSep := "KUBE-SEP-" + serviceName + "-" + fmt.Sprintf("%d", port) + "-" + fmt.Sprintf("%d", j)
		err = i.iptable.NewChain("nat", kubeSep)
		if err != nil {
			log.ErrorLog("KUBE-SEP NewChain error: " + err.Error())
			return err
		}

		// 在KUBE-SVC链中添加规则转发至KUBE-SEP链，默认采用随机算法
		if j == podsLen-1 {
			err = i.iptable.Append("nat", kubeSvc, "-j", kubeSep)
			if err != nil {
				log.ErrorLog("KUBE-SVC forward to KUBE_SEP error: " + err.Error())
				return err
			}
		} else {
			Probability := 1.0 / float64(podsLen-j)
			err = i.iptable.Append("nat", kubeSvc, "-m", "statistic", "--mode", "random", "--probability", fmt.Sprintf("%f", Probability), "-j", kubeSep)
			if err != nil {
				log.ErrorLog("KUBE-SVC forward to KUBE_SEP error: " + err.Error())
				return err
			}
		}

		// 在KUBE-SEP链中添加规则，转发至对应的pod
		err = i.iptable.Append("nat", kubeSep, "-d", clusterIP, "-p", string(protocol), "--dport", fmt.Sprintf("%d", port), "-j", "DNAT", "--to-destination", podIP+":"+fmt.Sprintf("%d", targetPort))
		if err != nil {
			log.ErrorLog("KUBE-SEP DNAT error: " + err.Error())
			return err
		}

		// 在KUBE-SEP链中添加规则，标记数据包
		err = i.iptable.Insert("nat", kubeSep, 1, "-d", clusterIP, "-p", string(protocol), "--dport", fmt.Sprintf("%d", port), "-j", "KUBE-MARK-MASQ")
		if err != nil {
			log.ErrorLog("KUBE-SEP Mark error: " + err.Error())
			return err
		}
	}

	return nil
}

func (i *iptableManager) UpdateService(updateEvent *entity.ServiceEvent) error {

	// 删除旧的service
	err := i.DeleteService(updateEvent)
	if err != nil {
		log.ErrorLog("UpdateService error: " + err.Error())
		return err
	}

	// 创建新的service
	err = i.CreateService(updateEvent)
	if err != nil {
		log.ErrorLog("UpdateService error: " + err.Error())
		return err
	}
	return nil
}

func (i *iptableManager) DeleteService(deleteEvent *entity.ServiceEvent) error {
	serviceName := deleteEvent.Service.Metadata.Name

	if deleteEvent.Service.Spec.Type == "ClusterIP" || deleteEvent.Service.Spec.Type == "" {
		for _, portConfig := range deleteEvent.Service.Spec.Ports {
			port := portConfig.Port
			// 在KUBE-SERVICES链中删除转发转发规则
			err := i.iptable.Delete("nat", "KUBE-SERVICES", "-d", deleteEvent.Service.Spec.ClusterIP, "-p", string(portConfig.Protocol), "--dport", fmt.Sprintf("%d", port), "-j", "KUBE-SVC-"+serviceName+"-"+fmt.Sprintf("%d", port))
			if err != nil {
				log.ErrorLog("DeleteService error: " + err.Error())
				return err
			}

			// 删除KUBE-SVC链
			kubeSvc := "KUBE-SVC-" + serviceName + "-" + fmt.Sprintf("%d", port)
			err = i.iptable.DeleteChain("nat", kubeSvc)
			if err != nil {
				log.ErrorLog("KUBE-SVC DeleteChain error: " + err.Error())
				return err
			}

			// 删除KUBE-SEP链
			for j := 0; j < len(i.service2Endpoint[serviceName]); j++ {
				err = i.iptable.DeleteChain("nat", "KUBE-SEP-"+serviceName+"-"+fmt.Sprintf("%d", port)+"-"+fmt.Sprintf("%d", j))
				if err != nil {
					log.ErrorLog("KUBE-SEP DeleteChain error: " + err.Error())
					return err
				}
			}

		}
	} else if deleteEvent.Service.Spec.Type == "NodePort" {
		for _, portConfig := range deleteEvent.Service.Spec.Ports {
			nodePort := portConfig.NodePort
			// 在KUBE-NODEPORTS链中删除转发规则
			err := i.iptable.Delete("nat", "KUBE-NODEPORTS", "-p", string(portConfig.Protocol), "--dport", fmt.Sprintf("%d", nodePort), "-j", "KUBE-SVC-"+serviceName+"-"+fmt.Sprintf("%d", portConfig.Port))
			if err != nil {
				log.ErrorLog("DeleteService error: " + err.Error())
				return err
			}

			// 删除KUBE-SVC链
			kubeSvc := "KUBE-SVC-" + serviceName + "-" + fmt.Sprintf("%d", portConfig.Port)
			err = i.iptable.DeleteChain("nat", kubeSvc)
			if err != nil {
				log.ErrorLog("KUBE-SVC DeleteChain error: " + err.Error())
				return err
			}

			// 删除KUBE-SEP链
			for j := 0; j < len(i.service2Endpoint[serviceName]); j++ {
				err = i.iptable.DeleteChain("nat", "KUBE-SEP-"+serviceName+"-"+fmt.Sprintf("%d", portConfig.Port)+"-"+fmt.Sprintf("%d", j))
				if err != nil {
					log.ErrorLog("KUBE-SEP DeleteChain error: " + err.Error())
					return err
				}
			}
		}
	} else {
		log.ErrorLog("DeleteService error: service type error")
		return fmt.Errorf("DeleteService error: service type error")

	}
	return nil
}
