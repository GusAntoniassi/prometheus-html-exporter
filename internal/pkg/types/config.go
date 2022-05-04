package types

type ExporterConfig struct {
	Targets      []TargetConfig
	GlobalConfig GlobalConfig `yaml:"global_config"`
}

type GlobalConfig struct {
	MetricNamePrefix string `yaml:"metric_name_prefix"`
	Port             int    `required:"true"`
	HasConfigFile    bool
}

type TargetConfig struct {
	Address               string `required:"true"`
	DecimalPointSeparator string `yaml:"decimal_point_separator"`
	ThousandsSeparator    string `yaml:"thousands_separator"`
	Metrics               []MetricConfig
}

type MetricConfig struct {
	Name     string `required:"true"`
	Help     string
	Type     string
	Selector string `required:"true"`
	Labels   map[string]string
}
