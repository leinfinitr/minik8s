package nginx

import (
	"fmt"
	"minik8s/pkg/apiObject"
	"strconv"
)

const (
	listenPort  = 80
	ServiceName = "nginx-svc"
)

func TranslateConfig(dns apiObject.Dns) string {
	var config string
	config += "server {\n"
	config += "    listen " + strconv.Itoa(listenPort) + ";\n"
	config += "    server_name " + dns.Spec.Host + ";\n"
	locationStr := "\tlocation %s {\n\t\tproxy_pass http://%s:%s/;\n\t}\n"
	for _, p := range dns.Spec.Paths {
		path := p.SubPath
		if path[0]!='/'{
			path = "/"+path
		}
		svcIp := p.SvcIp
		svcPort := p.SvcPort
		config += fmt.Sprintf(locationStr, path, svcIp, svcPort)
	}

	config += "}\n"
	return config
}