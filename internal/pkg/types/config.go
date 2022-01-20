package types

type ExporterConfig struct {
	ScrapeConfig ScrapeConfig `yaml:"scrape_config"`
	GlobalConfig GlobalConfig `yaml:"global_config"`
}

type GlobalConfig struct {
	MetricNamePrefix string `yaml:"metric_name_prefix"`
	Port             int
}

type ScrapeConfig struct {
	Name                  string `yaml:",omitempty"`
	Address               string
	Selector              string
	DecimalPointSeparator string       `yaml:"decimal_point_separator"`
	ThousandsSeparator    string       `yaml:"thousands_separator"`
	MetricConfig          MetricConfig `yaml:"metric"`
}

type MetricConfig struct {
	Name   string
	Help   string
	Type   string
	Labels map[string]string
}
