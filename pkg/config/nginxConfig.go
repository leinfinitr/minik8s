package config
import(
	"os"
)
var NginxPodYamlPath = os.Getenv("MINIK8S_PATH") + "nginx/yaml/nginx_pod.yaml"
var NginxServiceYamlPath = os.Getenv("MINIK8S_PATH") + "util/nginx/yaml/nginx_service.yaml"
var NginxDnsYamlPath = os.Getenv("MINIK8S_PATH") + "nginx/yaml/nginx_dns.yaml"