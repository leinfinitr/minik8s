package apiObject

const (
	// prometheus在master节点的配置文件路径
	PrometheusConfigPath = "/prometheus/prometheus/prometheus.yml"
	// 每个节点的NodeExporter的输出端口
	NodeExporterPort = 9190
)

// PrometheusConfig represents the structure of the Prometheus config file
type PrometheusConfig struct {
	Global struct {
		ScrapeInterval     string `yaml:"scrape_interval"`
		EvaluationInterval string `yaml:"evaluation_interval"`
	} `yaml:"global"`
	Alerting struct {
		Alertmanagers []struct {
			StaticConfigs []struct {
				Targets []string `yaml:"targets"`
			} `yaml:"static_configs"`
		} `yaml:"alertmanagers"`
	} `yaml:"alerting"`
	RuleFiles     []string       `yaml:"rule_files"`
	ScrapeConfigs []ScrapeConfig `yaml:"scrape_configs"`
}

type ScrapeConfig struct {
	JobName       string         `yaml:"job_name"`
	StaticConfigs []StaticConfig `yaml:"static_configs"`
}

type StaticConfig struct {
	Targets []string          `yaml:"targets"`
	Labels  map[string]string `yaml:"labels,omitempty"`
}
