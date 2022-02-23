package types

type ExporterConfig struct {
	Targets      []TargetConfig
	GlobalConfig GlobalConfig `yaml:"global_config"`
}

type GlobalConfig struct {
	MetricNamePrefix string `yaml:"metric_name_prefix"`
	Port             int
}

type TargetConfig struct {
	Address               string
	DecimalPointSeparator string `yaml:"decimal_point_separator"`
	ThousandsSeparator    string `yaml:"thousands_separator"`
	Metrics               []MetricConfig
}

type MetricConfig struct {
	Name     string
	Help     string
	Type     string
	Selector string
	Labels   map[string]string
}
