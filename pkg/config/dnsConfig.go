package config

const (
	DNS_Label_Key      = "dns"
	DNS_Label_Value    = "nginx"
	DNS_PATH           = "/etc/hosts"
	NginxConfigPath    = "/etc/nginx/nginx.conf"
	NginxMntConfigPath = "/mnt/nginx.conf"
	LocalConfigPath    = "/pvclient/:namespace/:name/nginx.conf"
	LocalOriginalPath  = "/root/nginx.conf"
)
